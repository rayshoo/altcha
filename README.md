<img src="./logo-black.altcha.svg" alt="ALTCHA" width="240" />

# ALTCHA

A lightweight ALTCHA challenge/verify service built with Go and Echo. Provides simple API endpoints for challenge generation and solution verification, with an optional demo UI.

Go(Echo) 기반의 경량 ALTCHA 챌린지/검증 서비스입니다. 챌린지 생성 및 솔루션 검증을 위한 간단한 API 엔드포인트를 제공하며, 선택적으로 데모 UI를 포함합니다.

## Quick Start

```bash
[ -f .env ] || cp .env.example .env
docker compose up --build
```

- API: http://localhost:3000
- Dashboard: http://localhost:9000 (`admin`/`admin`)
- Demo (`DEMO=true`): http://localhost:8000

## Documentation

[English](./docs/en/configuration.md) | [한국어](./docs/ko/configuration.md)

## Client Integration

[English](./docs/en/client-integration.md) | [한국어](./docs/ko/client-integration.md)

## Widget Customization

[English](./docs/en/widget-customization.md) | [한국어](./docs/ko/widget-customization.md)

## Dashboard

[English](./docs/en/dashboard.md) | [한국어](./docs/ko/dashboard.md)

## Production Architecture

[English](./docs/en/production-architecture.md) | [한국어](./docs/ko/production-architecture.md)

## License

[MIT](./LICENSE) © rayshoo
