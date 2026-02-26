# 운영 아키텍처

## 인증 흐름

ALTCHA는 챌린지 발급과 검증이 분리된 구조입니다. 운영 환경에서는 각 단계의 호출 주체를 구분하는 것이 중요합니다.

```
[브라우저]                    [UI 백엔드]              [ALTCHA 서버]
    |                             |                          |
    |--- GET /challenge ---------------------------------->  |  (1) 프론트엔드 → ALTCHA
    |<-- challenge JSON ----------------------------------|  |
    |                             |                          |
    |   (PoW 풀이, 브라우저 CPU)   |                          |
    |                             |                          |
    |--- POST /login ----------->|                           |  (2) 프론트엔드 → UI 백엔드
    |    (form + altcha payload)  |                           |
    |                             |--- GET /verify --------->|  (3) 백엔드 → ALTCHA (서버 to 서버)
    |                             |<-- 202/417 --------------|
    |<-- login result ------------|                           |
```

### 1단계: 챌린지 발급 (프론트엔드 → ALTCHA)

브라우저의 ALTCHA 위젯이 `GET /challenge`를 직접 호출합니다. 위젯이 PoW를 브라우저에서 풀어야 하므로 프론트엔드 호출이 필수입니다.

- 퍼블릭 Ingress로 노출 필요
- CORS 설정 필요 (`CORS_ORIGIN` 환경변수)

### 2단계: 폼 제출 (프론트엔드 → UI 백엔드)

사용자가 폼을 제출하면 `altcha` 페이로드가 함께 전송됩니다. 이 시점까지는 아직 검증되지 않은 상태입니다.

### 3단계: 솔루션 검증 (백엔드 → ALTCHA, 서버 to 서버)

UI 백엔드가 `GET /verify?altcha=<payload>`를 서버사이드에서 호출합니다.

**프론트엔드가 아닌 백엔드에서 호출해야 하는 이유:**

- **보안** — `/verify`를 퍼블릭에 노출하면 공격자가 직접 호출하여 토큰 소진 공격이 가능
- **신뢰성** — 프론트엔드의 검증 결과는 DevTools로 조작 가능. 백엔드에서 직접 확인해야 신뢰할 수 있음
- **네트워크** — 내부 네트워크로 호출하면 퍼블릭 인터넷을 경유하지 않음

## Kubernetes 네트워크 구성

### 같은 클러스터 내

UI 백엔드에서 ClusterIP Service를 통해 직접 호출:

```
http://altcha.<namespace>.svc.cluster.local:3000/verify
```

### VPC가 분리된 환경 (dev/stg/prd)

ALTCHA를 공유 서비스 VPC에 배포하거나, 각 환경에 개별 배포할 수 있습니다.

**방법 1: 환경별 개별 배포 (권장)**

각 EKS 클러스터에 ALTCHA를 배포합니다. 네트워크 의존성이 없고 장애가 격리됩니다.

```
[dev VPC]                [stg VPC]                [prd VPC]
 ├─ EKS                   ├─ EKS                   ├─ EKS
 │  ├─ UI App             │  ├─ UI App             │  ├─ UI App
 │  ├─ ALTCHA             │  ├─ ALTCHA             │  ├─ ALTCHA
 │  └─ Redis/Valkey       │  └─ Redis/Valkey       │  └─ Redis/Valkey
```

**방법 2: 공유 서비스 VPC**

VPC Peering이나 Transit Gateway를 통해 공유 ALTCHA 서비스에 접근합니다.

```
[shared VPC]
 ├─ ALTCHA + Redis
 └─ Internal ALB ←── VPC Peering ←── [dev/stg/prd VPC의 UI 백엔드]
```

## TLS 필수 (모바일 브라우저)

ALTCHA 위젯은 Web Worker를 사용하여 브라우저에서 PoW를 풀이합니다. 모바일 브라우저(Chrome, Safari)는 **비보안 컨텍스트(HTTP)에서 blob URL Worker 생성을 차단**하므로, 반드시 HTTPS(TLS)를 통해 서비스해야 합니다.

- `localhost`는 브라우저가 secure context로 취급하여 HTTP에서도 동작
- `http://<IP>` 또는 `http://<도메인>`은 insecure context로 Worker가 차단됨
- Ingress에 TLS를 설정하면 HTTPS로 접근하므로 문제 없음

## Ingress 분리

API 서버와 데모 서버의 Ingress를 분리하는 것을 권장합니다.

- **API Ingress** (`altcha.example.com`) — `/challenge` 엔드포인트를 브라우저에 노출. `/verify`는 클러스터 내부 Service로만 호출
- **Demo Ingress** (`altcha-demo.example.com`) — 데모 페이지 테스트용. 운영에서는 제거

`/verify`와 `/health/*`는 퍼블릭 Ingress에 노출하지 않습니다. UI 백엔드는 클러스터 내부 Service 주소로 직접 호출합니다.

## 환경변수 설정 예시

```bash
# 운영 환경
SECRET=<충분히 긴 랜덤 문자열>
CORS_ORIGIN=https://app.example.com,https://login.example.com
STORE=redis
REDIS_URL=redis://valkey-cluster.xxxxx.apne2.cache.amazonaws.com:6379
REDIS_CLUSTER=true
RATE_LIMIT=20
```
