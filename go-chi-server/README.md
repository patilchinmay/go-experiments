# Description

This is a basic go server written with go-chi.

Contains:

Please see [README.md](../README.md#toc) from parent folder.

# Setup

```bash
go mod tidy
```

# Run

```bash
> go run github.com/patilchinmay/go-experiments/go-chi-server

#OR

> go run main.go

2023-04-28T16:23:34.225171+09:00 INF ENV=local service="my app"
2023-04-28T16:23:34.225261+09:00 INF Deploying webhook PORT=8080 service="my app"
2023-04-28T16:23:34.225277+09:00 INF Deploying webhook HOST=0.0.0.0 service="my app"
2023-04-28T16:23:34.225309+09:00 INF Listening Addr=0.0.0.0:8080 service="my app"
2023-04-28T16:23:45.34459+09:00 INF Response: 200 OK httpRequest={"proto":"HTTP/1.1","remoteIP":"[::1]:53249","requestID":"FA21110345/DaoW3C4YXS-000001","requestMethod":"GET","requestPath":"/health","requestURL":"http://localhost:8080/health"} httpResponse={"bytes":1,"elapsed":0.020584,"status":200} service="my app"
2023-04-28T16:23:57.017593+09:00 INF Pong httpRequest={"proto":"HTTP/1.1","remoteIP":"[::1]:53251","requestID":"FA21110345/DaoW3C4YXS-000002","requestMethod":"GET","requestPath":"/ping","requestURL":"http://localhost:8080/ping"} service="my app"
2023-04-28T16:23:57.017644+09:00 INF Response: 200 OK httpRequest={"proto":"HTTP/1.1","remoteIP":"[::1]:53251","requestID":"FA21110345/DaoW3C4YXS-000002","requestMethod":"GET","requestPath":"/ping","requestURL":"http://localhost:8080/ping"} httpResponse={"bytes":4,"elapsed":0.0575,"status":200} service="my app"

```

# Verify
```bash
❯ curl localhost:8080/health
.%

❯ curl localhost:8080/ping
Pong%
```
