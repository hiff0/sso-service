FROM golang:1.23-alpine as builder

RUN apk add --no-cache git gcc musl-dev sqlite-dev supervisor

WORKDIR /app

COPY . .
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

RUN go mod tidy

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o sso ./cmd/sso/main.go
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o migrator ./cmd/migrator

EXPOSE 44044
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
# ENTRYPOINT ["./sso", "--config=./config/prod.yml"]

# FROM alpine:edge

# WORKDIR /app

# COPY --from=builder app/config/prod.yml config/prod.yml
# COPY --from=builder app/sso sso

# EXPOSE 44044

# ENTRYPOINT ["app/sso" "--config=./config/prod.yml"]