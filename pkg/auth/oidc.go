package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"altcha/pkg/config"
)

type OIDCProvider struct {
	cfg          *config.Config
	sessions     *SessionStore
	discovery    *oidcDiscovery
	jwksCache    *jwksCache
	callbackURL  string
	stateMu      sync.RWMutex
	pendingState map[string]*authState
}

type oidcDiscovery struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	EndSessionEndpoint    string `json:"end_session_endpoint"`
	JWKSURI               string `json:"jwks_uri"`
	Issuer                string `json:"issuer"`
}

type authState struct {
	CodeVerifier string
	CreatedAt    time.Time
}

type jwksCache struct {
	mu   sync.RWMutex
	keys map[string]interface{} // kid -> *rsa.PublicKey
	uri  string
}

func NewOIDCProvider(cfg *config.Config) *OIDCProvider {
	p := &OIDCProvider{
		cfg:          cfg,
		sessions:     NewSessionStore(),
		pendingState: make(map[string]*authState),
		jwksCache:    &jwksCache{keys: make(map[string]interface{})},
	}
	return p
}

func (p *OIDCProvider) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("altcha_session")
			if err != nil || cookie.Value == "" {
				return p.redirectToLogin(c)
			}

			session := p.sessions.Get(cookie.Value)
			if session == nil {
				return p.redirectToLogin(c)
			}

			// Check token expiry, try refresh
			if time.Now().After(session.ExpiresAt) {
				if session.RefreshToken != "" {
					if err := p.refreshTokens(cookie.Value, session); err != nil {
						p.sessions.Delete(cookie.Value)
						return p.redirectToLogin(c)
					}
				} else {
					p.sessions.Delete(cookie.Value)
					return p.redirectToLogin(c)
				}
			}

			c.Set("user", session.User)
			return next(c)
		}
	}
}

func (p *OIDCProvider) RegisterRoutes(e *echo.Echo) {
	e.GET("/auth/callback", p.handleCallback)
	e.GET("/auth/logout", p.handleLogout)
}

func (p *OIDCProvider) ensureDiscovery() error {
	if p.discovery != nil {
		return nil
	}

	disco := &oidcDiscovery{}

	// Fetch from well-known endpoint
	wellKnown := strings.TrimSuffix(p.cfg.AuthIssuer, "/") + "/.well-known/openid-configuration"
	resp, err := http.Get(wellKnown)
	if err != nil {
		return fmt.Errorf("fetch OIDC discovery: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OIDC discovery returned %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(disco); err != nil {
		return fmt.Errorf("decode OIDC discovery: %w", err)
	}

	// Allow env overrides
	if p.cfg.AuthAuthorizationEndpoint != "" {
		disco.AuthorizationEndpoint = p.cfg.AuthAuthorizationEndpoint
	}
	if p.cfg.AuthTokenEndpoint != "" {
		disco.TokenEndpoint = p.cfg.AuthTokenEndpoint
	}
	if p.cfg.AuthEndSessionEndpoint != "" {
		disco.EndSessionEndpoint = p.cfg.AuthEndSessionEndpoint
	}
	if p.cfg.AuthJWKSURI != "" {
		disco.JWKSURI = p.cfg.AuthJWKSURI
	}

	p.discovery = disco
	p.jwksCache.uri = disco.JWKSURI
	return nil
}

func (p *OIDCProvider) redirectToLogin(c echo.Context) error {
	if err := p.ensureDiscovery(); err != nil {
		return c.String(http.StatusInternalServerError, "OIDC discovery failed: "+err.Error())
	}

	state := generateRandomString(32)
	codeVerifier := generateRandomString(64)

	p.stateMu.Lock()
	p.pendingState[state] = &authState{
		CodeVerifier: codeVerifier,
		CreatedAt:    time.Now(),
	}
	p.stateMu.Unlock()

	// Cleanup old states periodically
	go p.cleanupStates()

	codeChallenge := computeCodeChallenge(codeVerifier)

	scheme := "https"
	if c.Request().TLS == nil {
		if fwd := c.Request().Header.Get("X-Forwarded-Proto"); fwd != "" {
			scheme = fwd
		} else {
			scheme = "http"
		}
	}
	callbackURL := fmt.Sprintf("%s://%s/auth/callback", scheme, c.Request().Host)

	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {p.cfg.AuthClientID},
		"redirect_uri":          {callbackURL},
		"scope":                 {"openid profile"},
		"state":                 {state},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
	}

	return c.Redirect(http.StatusFound, p.discovery.AuthorizationEndpoint+"?"+params.Encode())
}

func (p *OIDCProvider) handleCallback(c echo.Context) error {
	if err := p.ensureDiscovery(); err != nil {
		return c.String(http.StatusInternalServerError, "OIDC discovery failed")
	}

	code := c.QueryParam("code")
	state := c.QueryParam("state")

	p.stateMu.Lock()
	as, ok := p.pendingState[state]
	if ok {
		delete(p.pendingState, state)
	}
	p.stateMu.Unlock()

	if !ok {
		return c.String(http.StatusBadRequest, "Invalid or expired state")
	}

	scheme := "https"
	if c.Request().TLS == nil {
		if fwd := c.Request().Header.Get("X-Forwarded-Proto"); fwd != "" {
			scheme = fwd
		} else {
			scheme = "http"
		}
	}
	callbackURL := fmt.Sprintf("%s://%s/auth/callback", scheme, c.Request().Host)

	// Exchange code for tokens
	formData := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {callbackURL},
		"client_id":     {p.cfg.AuthClientID},
		"code_verifier": {as.CodeVerifier},
	}
	if p.cfg.AuthClientSecret != "" {
		formData.Set("client_secret", p.cfg.AuthClientSecret)
	}

	resp, err := http.PostForm(p.discovery.TokenEndpoint, formData)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Token exchange failed: "+err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return c.String(http.StatusInternalServerError, "Token exchange failed: "+string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		IDToken      string `json:"id_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to decode token response")
	}

	// Parse and verify ID token
	userInfo, err := p.parseIDToken(tokenResp.IDToken)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to verify ID token: "+err.Error())
	}

	if !IsAuthorized(userInfo, p.cfg) {
		return c.String(http.StatusForbidden, "Access denied")
	}

	expiresIn := time.Duration(tokenResp.ExpiresIn) * time.Second
	if expiresIn == 0 {
		expiresIn = 5 * time.Minute
	}
	sessionID := p.sessions.Create(userInfo, tokenResp.AccessToken, tokenResp.RefreshToken, expiresIn)

	c.SetCookie(&http.Cookie{
		Name:     "altcha_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Request().TLS != nil || c.Request().Header.Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	return c.Redirect(http.StatusFound, "/")
}

func (p *OIDCProvider) handleLogout(c echo.Context) error {
	cookie, err := c.Cookie("altcha_session")
	if err == nil && cookie.Value != "" {
		p.sessions.Delete(cookie.Value)
	}

	c.SetCookie(&http.Cookie{
		Name:     "altcha_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	if p.discovery != nil && p.discovery.EndSessionEndpoint != "" {
		params := url.Values{
			"client_id":                {p.cfg.AuthClientID},
			"post_logout_redirect_uri": {fmt.Sprintf("https://%s/", c.Request().Host)},
		}
		return c.Redirect(http.StatusFound, p.discovery.EndSessionEndpoint+"?"+params.Encode())
	}

	return c.Redirect(http.StatusFound, "/")
}

func (p *OIDCProvider) parseIDToken(tokenString string) (*UserInfo, error) {
	// Parse without verification first to get kid
	unverified, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("parse unverified: %w", err)
	}

	kid, _ := unverified.Header["kid"].(string)
	if kid == "" {
		return nil, fmt.Errorf("missing kid in token header")
	}

	key, err := p.jwksCache.getKey(kid)
	if err != nil {
		return nil, fmt.Errorf("get signing key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	}, jwt.WithIssuer(p.cfg.AuthIssuer))
	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	user := &UserInfo{}
	if sub, ok := claims["preferred_username"].(string); ok {
		user.Username = sub
	} else if sub, ok := claims["sub"].(string); ok {
		user.Username = sub
	}

	// Extract roles from realm_access.roles
	if ra, ok := claims["realm_access"].(map[string]interface{}); ok {
		if roles, ok := ra["roles"].([]interface{}); ok {
			for _, r := range roles {
				if s, ok := r.(string); ok {
					user.Roles = append(user.Roles, s)
				}
			}
		}
	}

	// Extract groups
	if groups, ok := claims["groups"].([]interface{}); ok {
		for _, g := range groups {
			if s, ok := g.(string); ok {
				user.Groups = append(user.Groups, s)
			}
		}
	}

	return user, nil
}

func (p *OIDCProvider) refreshTokens(sessionID string, session *Session) error {
	if err := p.ensureDiscovery(); err != nil {
		return err
	}

	formData := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {session.RefreshToken},
		"client_id":     {p.cfg.AuthClientID},
	}
	if p.cfg.AuthClientSecret != "" {
		formData.Set("client_secret", p.cfg.AuthClientSecret)
	}

	resp, err := http.PostForm(p.discovery.TokenEndpoint, formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh failed with status %d", resp.StatusCode)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	expiresIn := time.Duration(tokenResp.ExpiresIn) * time.Second
	if expiresIn == 0 {
		expiresIn = 5 * time.Minute
	}
	p.sessions.UpdateTokens(sessionID, tokenResp.AccessToken, tokenResp.RefreshToken, expiresIn)
	return nil
}

func (p *OIDCProvider) cleanupStates() {
	p.stateMu.Lock()
	defer p.stateMu.Unlock()
	now := time.Now()
	for k, v := range p.pendingState {
		if now.Sub(v.CreatedAt) > 10*time.Minute {
			delete(p.pendingState, k)
		}
	}
}

// JWKS key fetching

func (jc *jwksCache) getKey(kid string) (interface{}, error) {
	jc.mu.RLock()
	key, ok := jc.keys[kid]
	jc.mu.RUnlock()
	if ok {
		return key, nil
	}

	// Fetch and try again
	if err := jc.fetchKeys(); err != nil {
		return nil, err
	}

	jc.mu.RLock()
	key, ok = jc.keys[kid]
	jc.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("key %s not found in JWKS", kid)
	}
	return key, nil
}

type jwksResponse struct {
	Keys []jwkKey `json:"keys"`
}

type jwkKey struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func (jc *jwksCache) fetchKeys() error {
	resp, err := http.Get(jc.uri)
	if err != nil {
		return fmt.Errorf("fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("decode JWKS: %w", err)
	}

	jc.mu.Lock()
	defer jc.mu.Unlock()

	for _, k := range jwks.Keys {
		if k.Kty != "RSA" {
			continue
		}
		pubKey, err := parseRSAPublicKey(k.N, k.E)
		if err != nil {
			continue
		}
		jc.keys[k.Kid] = pubKey
	}
	return nil
}

func parseRSAPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:length]
}

func computeCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}
