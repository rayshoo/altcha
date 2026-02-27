# 대시보드

ALTCHA 대시보드는 `/challenge`, `/verify` 요청 통계를 시각화하는 웹 UI입니다. API 서버와 별도의 바이너리(`/dashboard`)로 실행되며, 동일한 PostgreSQL 데이터베이스를 공유합니다.

## 아키텍처

```
[브라우저] ──→ [Dashboard :9000] ──→ [PostgreSQL]
                                          ↑
[브라우저] ──→ [API Server :3000] ────────┘ (이벤트 기록)
```

- **API 서버**: `POSTGRES_URL` 설정 시 `/challenge`, `/verify` 요청을 PostgreSQL `events` 테이블에 기록
- **대시보드**: 동일한 PostgreSQL에서 통계를 조회하여 시각화

## 빠른 시작

```bash
docker compose up --build
```

- API: http://localhost:3000
- Dashboard: http://localhost:9000 (기본 인증: `admin`/`admin`)

## 환경변수

### Analytics (API 서버)

| 변수 | 필수 | 기본값 | 설명 |
|---|---|---|---|
| POSTGRES_URL | | | PostgreSQL 연결 URL. 설정 시 분석 활성화 |
| GEOIP_DB | | | GeoLite2-Country.mmdb 파일 경로. 국가별 통계 활성화 |

### Dashboard

| 변수 | 필수 | 기본값 | 설명 |
|---|---|---|---|
| POSTGRES_URL | O | | PostgreSQL 연결 URL |
| DASHBOARD_PORT | | `9000` | 대시보드 서버 포트 |
| AUTH_PROVIDER | O | | 인증 방식: `basic` 또는 `keycloak` |

### Basic 인증

| 변수 | 설명 |
|---|---|
| AUTH_USERNAME | Basic 인증 사용자명 |
| AUTH_PASSWORD | Basic 인증 비밀번호 |

### Keycloak (OIDC) 인증

| 변수 | 필수 | 설명 |
|---|---|---|
| AUTH_ISSUER | O | Keycloak realm URL (예: `https://keycloak.example.com/realms/myrealm`) |
| AUTH_CLIENT_ID | O | OIDC 클라이언트 ID |
| AUTH_CLIENT_SECRET | | 클라이언트 시크릿 (public 클라이언트는 불필요) |
| AUTH_PKCE | | PKCE 사용 여부 (기본: `true`) |
| AUTH_AUTHORIZATION_ENDPOINT | | 인가 엔드포인트 오버라이드 |
| AUTH_TOKEN_ENDPOINT | | 토큰 엔드포인트 오버라이드 |
| AUTH_END_SESSION_ENDPOINT | | 로그아웃 엔드포인트 오버라이드 |
| AUTH_JWKS_URI | | JWKS URI 오버라이드 |

### 접근 제어

| 변수 | 설명 |
|---|---|
| AUTH_ALLOWED_USERS | 허용할 사용자 (쉼표 구분) |
| AUTH_ALLOWED_GROUPS | 허용할 그룹 (쉼표 구분) |
| AUTH_ALLOWED_ROLES | 허용할 역할 (쉼표 구분) |

모두 미설정 시 인증된 모든 사용자에게 접근을 허용합니다. 하나라도 설정하면 OR 로직으로 검사합니다.

## GeoIP 설정

국가별 통계를 사용하려면 MaxMind GeoLite2-Country 데이터베이스가 필요합니다. 선택 사항이며, 미설정 시 대시보드의 Locations 섹션이 숨겨지고 나머지 기능은 정상 동작합니다.

1. [MaxMind](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data)에서 무료 계정 등록 (라이선스 키 발급 필요)
2. GeoLite2-Country.mmdb 다운로드
3. `GEOIP_DB` 환경변수에 파일 경로 설정

```env
GEOIP_DB=GeoLite2-Country.mmdb
```

> mmdb 파일(~6MB)은 바이너리이며 라이선스 상 재배포가 제한될 수 있으므로 git에 포함하지 않습니다 (`*.mmdb`가 `.gitignore`에 등록되어 있습니다). Docker/K8s 환경에서는 볼륨 마운트로 제공하세요.

## API 엔드포인트

대시보드는 내부적으로 다음 API를 제공합니다 (인증 필요):

- `GET /api/summary?from=YYYY-MM-DD&to=YYYY-MM-DD` — KPI 요약
- `GET /api/timeseries?from=YYYY-MM-DD&to=YYYY-MM-DD` — 일별 추이
- `GET /api/locations?from=YYYY-MM-DD&to=YYYY-MM-DD` — 국가/대륙별 통계

## Kubernetes 배포

```yaml
# dashboard-deploy.yaml에서 command로 /dashboard 실행
command: ["/dashboard"]
```

대시보드는 내부 네트워크에서만 접근 가능하도록 `nginx-internal` Ingress 클래스를 사용합니다.

## 대시보드 기능

- **KPI 카드**: Challenges, Verified, Failed, Avg Latency, 4XX Errors, 5XX Errors, Total Requests
- **추이 차트**: 일별 요청 수(막대) + 평균 지연 시간(선) 혼합 차트
- **위치 통계**: 대륙/국가별 요청 비율 (GeoIP 설정 시)
- **날짜 범위**: 7일/30일/90일/커스텀 선택
