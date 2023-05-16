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
   - [x] Separation of App and Server. This is a good practice and makes testing easier.
   - [x] Creation of App and Server using `Singleton` pattern.
   - [x] `Builder` pattern for App and Server setup with different methods such as `WithLogger, WithHost, WithPort etc.`.
   - [x] Implicit route registration.
     - Uses `Subrouter (go-chi-server/app/subrouters.go)` interface.
     - Registers the subrouter using side-effects (blank identifier import) in `main.go`.
     - `[]App.Subrouters` maintains a list of all `Subrouters` and mounts them using `App.MountSubrouters()`.
     - E.g. `go-chi-server/app/ping` package implements the `Subrouter` interface.
     - It has an `init()` for side-effect driven registration.

  **Traceability:**

   - [x] HTTP Request Logging (`httplog`).
   - [x] Configurable app logging (`zerolog`).
   - [x] End-to-end unique request id.
     - If the incoming request contains non-empty `X-Request-Id` header with value, it will be used.
     - Otherwise a unique id will be created using go-chi [RequestID](https://github.com/go-chi/chi/blob/master/middleware/request_id.go) middleware.
     - [RequestID](https://github.com/go-chi/chi/blob/master/middleware/request_id.go) is automatically set by `httplog.RequestLogger` in `go-chi-server/app/app.go:SetupMiddlewares()`.
     - Example in `GET /ping`.

  **Configuration:**

   - [x] Loads environment variables from `.env` file (`godotenv`).
   - [x] Injectable config from env vars for for `httplog`.
     - Structured json logging using `JSONLOGS` env var.
     - Log level setting using `LOGLEVEL` env var.
   - [x] Overrides the server (`/go-chi-server/server/server.go:Server`) config, sets defaults using env vars ([go-envconfig](https://github.com/sethvargo/go-envconfig)).

  **Docker:**

   - [x] Includes `Dockerfile`, `docker-compose.yaml` and `.dockerignore`.
   - [x] Uses `multi-stage` builds to reduce the size of resulting image.
   - [x] Leverages docker caching.
     - Copies `go.*` files first (e.g. `go.mod`, `go.sum`).
     - Then downloads the dependencies and caches them with `go mod download`.
     - Then copies the rest of the application files.
     - This makes sure that we reuse the cache for dependencies layer (even if there is a change in the application code).
   - [x] Uses `nonroot` user for running application.
   - [x] Uses `distroless` image for running application.

  **Tests and Docs:**

   - [x] Tests, coverage and how to run them (`ginkgo, gomega`).
   - [x] Explanatory comments and `godoc`.
   - [x] Code coverage.
     - To check code coverage, run `go tool cover -html=coverprofile.out`

  **Miscellaneous**:

   - [x] Graceful Shutdown / OS Interrupt signal handling in `main.go`.

3. [https-serving](./https-serving)
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
