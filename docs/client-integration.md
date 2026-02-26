# 클라이언트 통합

## 위젯 설치

폼에 ALTCHA 위젯을 추가하고 `challengeurl`을 이 서비스로 설정하세요.

```html
<script async defer src="https://cdn.jsdelivr.net/gh/altcha-org/altcha@main/dist/altcha.min.js" type="module"></script>
<form action="/your-submit" method="POST">
  <input name="email" placeholder="Email" />
  <altcha-widget challengeurl="http://localhost:3000/challenge"></altcha-widget>
  <button>Submit</button>
</form>
```

제출 시 요청 본문에 `altcha` 필드 값이 포함됩니다. 서버에서 `GET /verify?altcha=...`를 호출하고:

- `202 Accepted` → 검증 성공
- `417 Expectation Failed` → 유효하지 않거나 재사용된 토큰

## 수동 테스트

PowerShell:

```powershell
# $payload에 클라이언트의 altcha 값이 포함되어 있다고 가정
curl "http://localhost:3000/verify?altcha=$([uri]::EscapeDataString($payload))" -Method GET -UseBasicParsing
```

Unix/macOS:

```bash
curl -G \
  --data-urlencode "altcha=$payload" \
  http://localhost:3000/verify -i
```

## 위젯 스타일 커스터마이징

위젯의 모양을 변경하려면 [위젯 커스터마이징](./widget-customization.md) 문서를 참고하세요.
