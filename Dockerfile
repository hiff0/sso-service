FROM golang:1.23-alpine as builder

RUN apk add --no-cache git gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o sso ./cmd/sso/main.go
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o migrator ./cmd/migrator

FROM alpine:latest

RUN apk add --no-cache sqlite-dev supervisor

WORKDIR /app

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

COPY --from=builder app/sso sso
COPY --from=builder app/migrator migrator

COPY migrations migrations
COPY storage storage
COPY config/prod.yml config/prod.yml

RUN chmod +x sso migrator

EXPOSE 44044
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]