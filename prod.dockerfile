FROM golang:1.24.3-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -o main .

FROM debian:bookworm-slim AS production

RUN apt-get update && apt-get install -y ca-certificates bash && rm -rf /var/lib/apt/lists/*

RUN addgroup --system appgroup && adduser --system --ingroup appgroup appuser

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3001

CMD ["./main"]