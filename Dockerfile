FROM golang:1.24.3-alpine3.20 AS dev

RUN apk add --no-cache git

RUN git config --global --add safe.directory /app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0

EXPOSE 3001

ENTRYPOINT ["tail", "-f", "/dev/null"]

FROM golang:1.24.3-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -trimpath -ldflags="-s -w" -o main .

FROM debian:bookworm-slim AS prod

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

RUN addgroup --system appgroup && adduser --system --ingroup appgroup appuser

WORKDIR /app

COPY --from=builder /app/main .

USER appuser

EXPOSE 3001

CMD ["./main"]
