# 설정

## 환경변수

| 변수 | 필수 | 기본값 | 설명 |
|---|---|---|---|
| SECRET | O | `$ecret.key` | HMAC 서명/검증 키. 프로덕션에서는 반드시 변경 |
| ALGORITHM | | `SHA-256` | 해시 알고리즘: `SHA-256`, `SHA-512`, `SHA-1` |
| PORT | | `3000` | API 서버 포트 |
| EXPIREMINUTES | | `10` | 챌린지 만료 시간(분). 사용자가 이 시간 안에 제출해야 함 |
| COMPLEXITY | | `1000000` | PoW 난이도. 클수록 클라이언트 브라우저 연산 시간 증가 |
| MAXRECORDS | | `1000` | 토큰 재사용 방지 캐시 크기 (memory/sqlite만 해당, redis는 TTL 사용) |
| CORS_ORIGIN | | `*` | 허용할 오리진 (쉼표 구분) |
| RATE_LIMIT | | `0` (무제한) | IP당 초당 요청 수 제한 |
| STORE | | `memory` | 토큰 저장소: `memory`, `sqlite`, `redis` |
| SQLITE_PATH | | `data/altcha.db` | SQLite 파일 경로 (STORE=sqlite 시) |
| REDIS_URL | | `redis://localhost:6379` | Redis 연결 URL (STORE=redis 시) |
| REDIS_CLUSTER | | `false` | `true` 시 클러스터 모드 사용 (ElastiCache, Valkey 등) |
| LOG_LEVEL | | `info` | `info`: API 로그만, `debug`: API + 데모 로그 |
| DEMO | | `false` | `true` 시 포트 8000에서 데모 UI 시작 |
| POSTGRES_URL | | | PostgreSQL 연결 URL. 설정 시 분석(analytics) 활성화 |
| GEOIP_DB | | | GeoLite2-Country.mmdb 경로 (선택, 국가별 통계 활성화) |
| DASHBOARD_PORT | | `9000` | 대시보드 서버 포트 |
| AUTH_PROVIDER | | | 대시보드 인증 방식: `basic` 또는 `keycloak` |
| AUTH_USERNAME | | | Basic 인증 사용자명 |
| AUTH_PASSWORD | | | Basic 인증 비밀번호 |
| AUTH_ISSUER | | | Keycloak realm URL |
| AUTH_CLIENT_ID | | | OIDC 클라이언트 ID |
| AUTH_CLIENT_SECRET | | | OIDC 클라이언트 시크릿 |
| AUTH_PKCE | | `true` | PKCE 사용 여부 |

### 챌린지 관련

- **EXPIREMINUTES**: `/challenge`에서 생성한 토큰의 유효 시간(분). 기본 10분이면 사용자가 챌린지를 받고 10분 안에 폼을 제출해야 합니다. 만료된 토큰은 `/verify`에서 자동으로 거부됩니다.
- **COMPLEXITY**: PoW 난이도를 조절하는 최대 숫자. 클라이언트 브라우저가 0부터 이 숫자 사이에서 정답을 찾아야 합니다. 값이 클수록 풀이 시간이 길어져 봇 공격 비용이 올라가지만, 일반 사용자 체감 지연도 증가합니다.
- **MAXRECORDS**: 검증된 토큰을 몇 개까지 기억할지 설정합니다. 캐시가 가득 차면 가장 오래된 것부터 삭제(FIFO). Redis 스토어에서는 TTL로 자동 만료하므로 이 값을 무시합니다.

## 환경변수 제공 방법

- `.env` 파일 (프로젝트 루트)
- `compose.yaml`의 `environment` 섹션
- 셸 환경변수 직접 설정

예시 `.env`:

```env
SECRET=change-me-to-a-long-random-string
PORT=3000
EXPIREMINUTES=10
COMPLEXITY=1000000
MAXRECORDS=1000
STORE=memory
DEMO=false
```

## 토큰 저장소

토큰 재사용 방지를 위한 세 가지 저장소 백엔드를 제공합니다.

### memory (기본)

인메모리 FIFO 캐시. 가장 단순하며 외부 의존성이 없습니다.

- 파드 재시작 시 캐시가 초기화됩니다.
- 단일 인스턴스 환경에 적합합니다.

```env
STORE=memory
MAXRECORDS=1000
```

### sqlite

파일 기반 영속 저장소. 파드가 재시작되거나 이동해도 PV(Persistent Volume)를 마운트하면 데이터가 유지됩니다.

- 단일 인스턴스 + 영속성이 필요할 때 사용합니다.
- 순수 Go 드라이버(CGO 불필요)를 사용합니다.

```env
STORE=sqlite
SQLITE_PATH=data/altcha.db
MAXRECORDS=1000
```

### redis

공유 저장소. 여러 인스턴스(파드)가 동일한 Redis를 바라보므로 수평 확장이 가능합니다.

- `EXPIREMINUTES`를 TTL로 사용하여 자동 만료됩니다.
- `MAXRECORDS` 설정은 무시됩니다 (TTL이 정리를 담당).

단일 노드:

```env
STORE=redis
REDIS_URL=redis://redis-host:6379
```

클러스터 (ElastiCache, Valkey 등 단일 엔드포인트):

```env
STORE=redis
REDIS_URL=redis://cluster-endpoint:6379
REDIS_CLUSTER=true
```

클러스터 (여러 노드 직접 지정, 자동 감지):

```env
STORE=redis
REDIS_URL=redis://node1:6379,redis://node2:6379,redis://node3:6379
```
