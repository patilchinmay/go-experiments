# Description

[go-chi-server](../go-chi-server) + HTTPS serving

Contains:

Please see [README.md](../README.md#toc) from parent folder.

# Key and Certificate creation

Ref:

```bash
cd certs

# Generate private key
openssl genrsa -out server.key 2048

# Generate public key (certificate) using the private key
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

# Setup

```bash
go mod tidy
```

# Run

```bash
❯ make run
go run main.go --tlscert certs/server.crt --tlskey certs/server.key
2023-05-01T21:21:26.297724+09:00 INF ENV=local service="my app"
2023-05-01T21:21:26.298737+09:00 INF Listening Addr=0.0.0.0:443 service="my app"

```

# Verify
```bash
❯ curl https://localhost:443/health --insecure
.%

❯ curl https://localhost:443/ping --insecure
Pong%
```

# Testing

## required utilities for testing

Install `godotenv, ginkgo and gomega`
https://onsi.github.io/ginkgo/

```bash
go install github.com/joho/godotenv/cmd/godotenv
go install github.com/onsi/ginkgo/v2/ginkgo
go get github.com/onsi/gomega/...
```

## run tests

`make test`

OR

`godotenv -f ./.env ginkgo -v -r --cover`

To check coverage via HTML: `go tool cover -html=coverprofile.out`

# References

- https://github.com/denji/golang-tls
- https://medium.com/jspoint/a-brief-overview-of-the-tcp-ip-model-ssl-tls-https-protocols-and-ssl-certificates-d5a6269fe29e
