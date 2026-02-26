# Production Architecture

## Authentication Flow

ALTCHA separates challenge issuance from verification. In production, it's important to distinguish who calls each endpoint.

```
[Browser]                    [UI Backend]              [ALTCHA Server]
    |                             |                          |
    |--- GET /challenge ---------------------------------->  |  (1) Frontend → ALTCHA
    |<-- challenge JSON ----------------------------------|  |
    |                             |                          |
    |   (PoW solving, browser CPU)|                          |
    |                             |                          |
    |--- POST /login ----------->|                           |  (2) Frontend → UI Backend
    |    (form + altcha payload)  |                           |
    |                             |--- GET /verify --------->|  (3) Backend → ALTCHA (server-to-server)
    |                             |<-- 202/417 --------------|
    |<-- login result ------------|                           |
```

### Step 1: Challenge Issuance (Frontend → ALTCHA)

The ALTCHA widget in the browser calls `GET /challenge` directly. This must be a frontend call since the widget solves PoW in the browser.

- Requires public Ingress exposure
- Requires CORS configuration (`CORS_ORIGIN` env var)

### Step 2: Form Submission (Frontend → UI Backend)

When the user submits the form, the `altcha` payload is sent along with it. At this point, the solution has not been verified yet.

### Step 3: Solution Verification (Backend → ALTCHA, server-to-server)

The UI backend calls `GET /verify?altcha=<payload>` server-side.

**Why this should be a backend call, not a frontend call:**

- **Security** — Exposing `/verify` publicly allows attackers to exhaust tokens directly
- **Reliability** — Frontend verification results can be tampered with via DevTools. Only backend verification is trustworthy
- **Network** — Internal network calls avoid traversing the public internet

## Kubernetes Network Configuration

### Same Cluster

Call directly via ClusterIP Service from the UI backend:

```
http://altcha.<namespace>.svc.cluster.local:3000/verify
```

### Separate VPCs (dev/stg/prd)

ALTCHA can be deployed in a shared services VPC or individually per environment.

**Option 1: Per-environment deployment (recommended)**

Deploy ALTCHA to each EKS cluster. No cross-VPC network dependencies and failures are isolated.

```
[dev VPC]                [stg VPC]                [prd VPC]
 ├─ EKS                   ├─ EKS                   ├─ EKS
 │  ├─ UI App             │  ├─ UI App             │  ├─ UI App
 │  ├─ ALTCHA             │  ├─ ALTCHA             │  ├─ ALTCHA
 │  └─ Redis/Valkey       │  └─ Redis/Valkey       │  └─ Redis/Valkey
```

**Option 2: Shared Services VPC**

Access a shared ALTCHA service via VPC Peering or Transit Gateway.

```
[shared VPC]
 ├─ ALTCHA + Redis
 └─ Internal ALB ←── VPC Peering ←── [UI backends in dev/stg/prd VPCs]
```

## TLS Required (Mobile Browsers)

The ALTCHA widget uses Web Workers to solve PoW in the browser. Mobile browsers (Chrome, Safari) **block blob URL Worker creation in insecure contexts (HTTP)**, so TLS (HTTPS) is mandatory.

- `localhost` is treated as a secure context by browsers, so HTTP works locally
- `http://<IP>` or `http://<domain>` is an insecure context where Workers are blocked
- Ingress with TLS provides HTTPS access, resolving this issue

## Ingress Separation

It is recommended to separate API and demo Ingress resources.

- **API Ingress** (`altcha.example.com`) — Exposes `/challenge` endpoint to browsers. `/verify` should only be called via internal cluster Service
- **Demo Ingress** (`altcha-demo.example.com`) — For demo page testing. Remove in production

Do not expose `/verify` or `/health/*` through the public Ingress. The UI backend should call these via the internal cluster Service address.

## Environment Variable Example

```bash
# Production
SECRET=<long-random-string>
CORS_ORIGIN=https://app.example.com,https://login.example.com
STORE=redis
REDIS_URL=redis://valkey-cluster.xxxxx.apne2.cache.amazonaws.com:6379
REDIS_CLUSTER=true
RATE_LIMIT=20
```
