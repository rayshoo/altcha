<img src="./logo-black.altcha.svg" alt="ALTCHA" width="240" />

# ALTCHA

Go(Echo) 기반의 경량 ALTCHA 챌린지/검증 서비스입니다. 챌린지 생성 및 솔루션 검증을 위한 간단한 API 엔드포인트를 제공하며, 선택적으로 데모 UI를 포함합니다.

- 런타임: Go
- 포트: 3000 (API), 8080 (선택적 데모)
- 라이브러리: [altcha-lib-go](https://github.com/altcha-org/altcha-lib-go)

## 빠른 시작

```powershell
Copy-Item .env.example .env -ErrorAction SilentlyContinue
docker compose up --build
```

```bash
# Unix/macOS
[ -f .env ] || cp .env.example .env
docker compose up --build
```

- API: http://localhost:3000
- 데모 (`DEMO=true` 시): http://localhost:8080

## 엔드포인트

| 메서드 | 경로 | 응답 | 설명 |
|---|---|---|---|
| GET | `/` | 204 | 라이브니스 프로브 |
| GET | `/health` | 200 JSON | 서버 상태, 버전, Go 런타임 정보 |
| GET | `/challenge` | 200 JSON | ALTCHA 챌린지 생성 |
| GET | `/verify?altcha=<payload>` | 202 / 417 | 솔루션 검증 (202: 성공, 417: 실패/재사용) |

## 설정

주요 환경변수:

| 변수 | 기본값 | 설명 |
|---|---|---|
| SECRET | `$ecret.key` | HMAC 키 (필수 변경) |
| PORT | `3000` | API 포트 |
| STORE | `memory` | 토큰 저장소: `memory`, `sqlite`, `redis` |
| RATE_LIMIT | `0` (무제한) | IP당 초당 요청 수 제한 |
| DEMO | `false` | 데모 UI 활성화 |

전체 설정 목록은 [설정 문서](./docs/configuration.md)를 참고하세요.

## 문서

- [설정 및 토큰 저장소](./docs/configuration.md) — 환경변수, memory/sqlite/redis 스토어
- [클라이언트 통합](./docs/client-integration.md) — 위젯 설치 및 검증 API 사용법
- [위젯 커스터마이징](./docs/widget-customization.md) — CSS 변수, HTML 속성

## 로컬 빌드

```bash
make build    # 빌드
make run      # 실행
make dev      # 개발 (air 핫리로드)
```

## 라이선스

MIT © rayshoo
