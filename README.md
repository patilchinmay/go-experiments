# go-experiments

This repository is a [go workspace](https://go.dev/doc/tutorial/workspaces) for ease of running the examples.

# How to read/understand this?

Basic familiarity of golang and backend development is assumed.

Each repository has a README.md with description, relevant details to setup and run.

# Table of Contents

1. [hello](./hello/)
   - [x] Verify that setup works with a hello world application.

2. [go-chi-server](./go-chi-server/)

   **Basics:**

   - [x] Basic Go-chi Server (`BCS`)
   - [x] Separation of App and Server as well as creation using `Singleton` pattern
   - [x] Builder pattern for App and Server creation with different methods such as `WithLogger, WithHost, WithPort etc.`

  **Traceability:**

   - [x] HTTP Request Logging (`httplog`)
   - [x] Configurable app logging (`zerolog`)
   - [x] End-to-end unique request id
     - If the incoming request contains `X-Request-Id` header with value, it will be used.
     - Otherwise a unique id will be created using go-chi [RequestID](https://github.com/go-chi/chi/blob/master/middleware/request_id.go) middleware.
     - [RequestID](https://github.com/go-chi/chi/blob/master/middleware/request_id.go) is automatically set by httplog.RequestLogger
     - Example in `GET /ping`

  **Configuration:**

   - [x] Injectable config from env vars for structured json logging and log level
   - [x] Loads environment variables from `.env` file (`godotenv`)
   - [x] Overrides the server (`/go-chi-server/server/server.go`) config, sets defaults using env vars ([go-envconfig](https://github.com/sethvargo/go-envconfig))


   - [x] Graceful Shutdown / OS Interrupt signal handling in `main`
   - [x] Implicit route registration using `Subrouter (./go-chi-server/app/subrouters.go)` interface in `./go-chi-server/app/routes/*`

   **Tests and Docs:**

   - [x] Tests, coverage and how to run them (`ginkgo, gomega`)
   - [x] Explanatory comments and `godoc`

1. [https-serving](./https-serving)
   - [x]  Base [go-chi-server](./go-chi-server/)
   - [x]  HTTPS serving
   - [x]  Automatic certificate [reloading](https://opensource.com/article/22/9/dynamically-update-tls-certificates-golang-server-no-downtime) on certificate changes (e.g. renewal)

# Possible Improvements

- [ ] Config loading
- [ ] Input parameter validation/sanitization
- [ ] Parameter Object
- [ ] Context parameter
- [ ] Middleware
- [ ] Database/ORM/Repository Layer
- [ ] Database disconnection/disruption
- [ ] JSON handling
- [ ] Route Versioning
- [ ] Cloud Native Golang Constructs e.g. retry, switch-breaker etc.
- [ ] Load Testing
- [ ] Code Quality / Static Analysis (e.g. Sonarqube, Codeclimate etc.)
- [ ] Sentry
- [ ] Feature Toggling
- [ ] UML diagram
- [ ] DB Connection with Connection pool
- [ ] Exposing metrics to Prometheus using go's runtime/metrics package.
- [ ] Timing out HTTP connections
- [ ] Handle DB disconnections
- [ ] Forcing HTTP 2.0 / 1.2
- [ ] Idle, Read and write timeout in http.Server
- [ ] Verify that http.ListenAndServe fields each new request on a separate goroutine by sending multiple requests and collecting runtime/metrics.
- [ ] Is it possible to print the goroutine ID? For simplifying above.
- [ ] Websocket server and client. Gorilla websocket library.
- [ ] DDoS and how to avoid it.
- [ ] Caching (local, redis etc)
