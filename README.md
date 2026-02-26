<img src="./logo-black.altcha.svg" alt="ALTCHA" width="240" />

# ALTCHA Docker

Go(Echo) 기반의 경량 ALTCHA 챌린지/검증 서비스입니다. 챌린지 생성 및 솔루션 검증을 위한 간단한 API 엔드포인트를 제공하며, 선택적으로 데모 UI를 포함합니다.

- 런타임: Go
- 포트: 3000 (API), 8080 (선택적 데모)
- 라이브러리: [altcha-lib-go](https://github.com/altcha-org/altcha-lib-go)

## 빠른 시작

Docker Compose 사용 (권장):

```powershell
# 선택적으로 .env 파일을 생성하거나 셸에서 변수를 설정하세요
Copy-Item .env.example .env -ErrorAction SilentlyContinue
# 스택 시작
docker compose up --build
```

```bash
# Unix/macOS (bash/zsh)
[ -f .env ] || cp .env.example .env
docker compose up --build
```

- API: http://localhost:3000
- 데모 (활성화 시): http://localhost:8080

시크릿을 임시로 변경하려면:

```powershell
$env:SECRET = "your-very-long-random-key"; docker compose up --build
```

```bash
SECRET="your-very-long-random-key" docker compose up --build
```

## 설정

서비스는 다음 환경변수를 사용합니다:

- SECRET (필수): 챌린지 서명/검증에 사용되는 HMAC 키. 기본값 `$ecret.key`는 프로덕션에서 사용하지 마세요.
- PORT: API 포트, 기본값 3000.
- EXPIREMINUTES: 챌린지 만료 시간(분), 기본값 10.
- MAXNUMBER: PoW 난이도를 위한 최대 숫자, 기본값 1000000.
- MAXRECORDS: 인메모리 일회용 토큰 캐시 크기, 기본값 1000.
- CORS_ORIGIN: 허용할 오리진 (쉼표 구분), 미설정 시 `*`.
- DEMO: `true`로 설정하면 포트 8080에서 데모 UI를 시작합니다.

`.env` 파일, `compose.yaml`의 environment 섹션, 또는 셸 환경변수로 제공할 수 있습니다.

예시 `.env`:

```env
SECRET=change-me-to-a-long-random-string
PORT=3000
EXPIREMINUTES=10
MAXNUMBER=1000000
MAXRECORDS=1000
DEMO=false
```

## 엔드포인트

- GET /

  - 204 No Content 반환. 라이브니스 프로브 엔드포인트.

- GET /health

  - 서버 상태, 버전, Go 런타임 정보를 JSON으로 반환합니다.
  - 200 OK.

- GET /challenge

  - altcha-lib으로 생성된 ALTCHA 챌린지 JSON을 반환합니다.
  - 200 OK와 챌린지 페이로드.

- GET /verify?altcha=\<payload\>
  - 제공된 ALTCHA 솔루션을 검증합니다.
  - 성공 시 202 Accepted.
  - 실패 또는 토큰 재사용 시 417 Expectation Failed (인메모리 캐시로 일회성 사용 강제).

참고:

- CORS는 기본적으로 열려 있습니다 (origin: \*). `CORS_ORIGIN`으로 특정 오리진을 지정할 수 있습니다.
- 레코드 재사용 방지는 인메모리에 저장됩니다. 스케일 아웃 또는 재시작 시 캐시가 초기화됩니다. 프로덕션에서는 공유 저장소와 함께 사용하세요.

## 데모 UI (선택)

`DEMO=true`를 설정하면 내장 데모 페이지를 활성화합니다. http://localhost:8080 에서 포트 3000의 API를 가리키는 ALTCHA 위젯이 포함된 간단한 HTML 폼을 제공합니다.

## 클라이언트 통합 예시

폼에 위젯을 추가하고 `challengeurl`을 이 서비스로 설정하세요:

```html
<script async defer src="https://cdn.jsdelivr.net/gh/altcha-org/altcha@main/dist/altcha.min.js" type="module"></script>
<form action="/your-submit" method="POST">
  <input name="email" placeholder="Email" />
  <altcha-widget challengeurl="http://localhost:3000/challenge"></altcha-widget>
  <button>Submit</button>
  <!-- 제출 시 요청 본문에 `altcha` 필드 값을 포함하세요 -->
  <!-- 서버에서 GET /verify?altcha=... 를 호출하고 202를 성공으로 처리하세요 -->
  <!-- 417은 유효하지 않거나 재사용된 토큰을 의미합니다 -->
</form>
```

PowerShell에서 검증을 수동으로 테스트:

```powershell
# $payload에 클라이언트의 `altcha` 값이 포함되어 있다고 가정
curl "http://localhost:3000/verify?altcha=$([uri]::EscapeDataString($payload))" -Method GET -UseBasicParsing
```

Unix/macOS에서 수동 테스트:

```bash
# $payload에 클라이언트의 `altcha` 값이 포함되어 있다고 가정
curl -G \
  --data-urlencode "altcha=$payload" \
  http://localhost:3000/verify -i
```

성공 시 202, 실패/재사용 시 417을 기대하세요.

## Docker 없이 빌드 및 실행

Go가 설치되어 있다면 로컬에서 실행 가능합니다:

```bash
make build
make run
```

개발 중 라이브 리로드:

```bash
make dev
```

## 프로덕션 참고사항

- SECRET을 강력하고 고유한 값으로 변경하세요. 기본값을 사용하지 마세요.
- 컨테이너 앞에 TLS 종료를 두고, 필요 시 /verify 접근을 제한하세요.
- 수평 확장 시 인메모리 토큰 캐시를 공유 저장소로 대체하세요.
- 이미지 버전을 고정하고, 여러 아키텍처에 배포하는 경우 멀티 아키텍처 빌드를 고려하세요.

## 라이선스

MIT © Umami Creative GmbH
