FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags "-X altcha/pkg/handler.Version=${VERSION}" -o /server ./cmd/server
RUN CGO_ENABLED=0 go build -ldflags "-X altcha/pkg/handler.Version=${VERSION}" -o /dashboard ./cmd/dashboard

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /server /server
COPY --from=build /dashboard /dashboard
COPY .env .env
COPY web/ web/
EXPOSE 3000
CMD ["/server"]
