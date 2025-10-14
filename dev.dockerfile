FROM golang:1.24.3-alpine3.20 AS builder

RUN apk update && apk add --no-cache git bash

RUN git config --global color.ui auto \
    && git config --global --add safe.directory /app

WORKDIR /app

COPY .git .git

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main .

EXPOSE 3001

ENTRYPOINT ["bash", "-c", "tail -f /dev/null"]