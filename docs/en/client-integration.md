# Client Integration

## Widget Installation

Add the ALTCHA widget to your form and point `challengeurl` to this service.

```html
<script async defer src="https://cdn.jsdelivr.net/gh/altcha-org/altcha@main/dist/altcha.min.js" type="module"></script>
<form action="/your-submit" method="POST">
  <input name="email" placeholder="Email" />
  <altcha-widget challengeurl="http://localhost:3000/challenge"></altcha-widget>
  <button>Submit</button>
</form>
```

On form submission, the `altcha` field value is included in the request body. Call `GET /verify?altcha=...` from your server:

- `202 Accepted` → Verification successful
- `417 Expectation Failed` → Invalid or reused token

## Manual Testing

PowerShell:

```powershell
# Assuming $payload contains the altcha value from the client
curl "http://localhost:3000/verify?altcha=$([uri]::EscapeDataString($payload))" -Method GET -UseBasicParsing
```

Unix/macOS:

```bash
curl -G \
  --data-urlencode "altcha=$payload" \
  http://localhost:3000/verify -i
```

## Widget Styling

To customize the widget appearance, see [Widget Customization](./widget-customization.md).
