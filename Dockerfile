FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
	CGO_ENABLED=0 go build -o /out/api ./cmd/api

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
	CGO_ENABLED=0 go build -o /out/migrate ./cmd/migrate

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /out/api /usr/local/bin/api
COPY --from=builder /out/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

ENTRYPOINT ["api"]