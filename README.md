<img src="./logo-black.altcha.svg" alt="ALTCHA" width="240" />

# ALTCHA

A self-hosted ALTCHA proof-of-work challenge server built with Go. Provides challenge generation and solution verification APIs, with an analytics dashboard (PostgreSQL, GeoIP) and Keycloak(OIDC) authentication support.

Go 기반의 셀프호스팅 ALTCHA 챌린지 서버입니다. 챌린지 생성 및 솔루션 검증 API를 제공하며, 분석 대시보드(PostgreSQL, GeoIP)와 Keycloak(OIDC) 인증을 지원합니다.

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
