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
   - [x] Implicit route registration.
     - Uses `Subrouter (go-chi-server/app/subrouters.go)` struct.
     - Registers the subrouter using side-effects (blank identifier import) in `main.go`.
     - `[]App.Subrouters` maintains a list of all `Subrouters` and mounts them using `App.MountSubrouters()`.
     - E.g. `go-chi-server/app/ping` package creates and configures its own `Subrouter` in the `init()` function.
     - This is a side-effect driven registration for subrouter onto main router.

   **Software Patterns:**

   - [x] **`Singleton`**: Creation of App and Server.
   - [x] **`Builder`**: App and Server setup with different methods such as `WithLogger, WithHost, WithPort etc.`.

   **Cloud Native Patterns:**

   - [x] **`Retry`**: Defined in [cloudnativepatterns](./cloudnativepatterns/) and used in [UserService](./go-chi-server/app/user/service.go).
   - [ ] **`Circuit Breaker`**:

   **Traceability:**

   - [x] HTTP Request Logging (`httplog`).
   - [x] Configurable app logging (`zerolog`).
   - [x] End-to-end unique request id.
     - [x] **Middlewares:**
     - If the incoming request contains non-empty `X-Request-Id` header with value, it will be used.
     - Otherwise a unique id will be created using go-chi [RequestID](https://github.com/go-chi/chi/blob/master/middleware/request_id.go) middleware.
     - [RequestID](https://github.com/go-chi/chi/blob/master/middleware/request_id.go) is automatically set by `httplog.RequestLogger` in `go-chi-server/app/app.go:SetupMiddlewares()`.
     - All responses contain `Request-Id` header. This header is added using the custom middleware `go-chi-server/app/middlewares/requestid.go`.

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
   - [x] Database
     - [x] Dockerized `Postgres` and `Pgadmin`.
     - [x] Healthcheck for `Postgres`.
     - [x] `depends_on` check for dependent services.

   **Database:**

   - [x] Service + Repository Layer (in [user](go-chi-server/app/user) module)
   - [x] Example in [user](go-chi-server/app/user) module with a CRUD REST API
   - [x] Database ORM ([gorm](https://github.com/go-gorm/gorm))
   - [x] DB Connection with Connection pool
   - [ ] Migrations ([golang-migrate](https://github.com/golang-migrate/migrate))
   - [ ] Database disconnection/disruption
   - [ ] Atomic Transactions

   **Tests, Benchmark and Docs:**

   - [x] Tests, coverage (`ginkgo, gomega`) and how to run them in [Makefile](go-chi-server/Makefile).
   - [x] SQL mock testing with [go-sqlmock](go-chi-server/app/user/repository_test.go)
   - [x] Testing with [SQLite](go-chi-server/app/user/repository_sqllite_test.go)
   - [x] Mock generation and testing using [gomock/mockgen](go-chi-server/app/user/service_test.go)
   - [x] HTTP handlers testing with [httptest](go-chi-server/app/user/handlers_test.go)
   - [x] Speed up testing using [mock clock](https://github.com/benbjohnson/clock) for retries in [UserService](./go-chi-server/app/user/service.go) so we don't actually wait for retry time intervals.
   - [x] Explanatory comments and `godoc`.
   - [x] **Code Coverage**
     - Run tests: `make test`
     - To check code coverage: `go tool cover -html=coverprofile.out`
   - [x] **Load Testing** (`apache benchmark`)
     - Start app: `make run`
     - Run `ab`: `docker run --rm jordi/ab -s 5 -k -c 500 -n 10000 http://host.docker.internal:8080/ping`
     - Why `ab` in docker? Because of [issues](https://serverfault.com/questions/806585/why-is-ab-erroring-with-apr-socket-recv-connection-reset-by-peer-54-on-osx) with ab in macos.
     - Time per request (mean):
       - Time per concurrent block `(-c)`
       - A.K.A. Average time taken per `(-c)` requests
     - Time per request (mean, across all concurrent requests):
       - Total Time / Total No of Requests
       - A.K.A. Average time taken for single request

   **Miscellaneous**:

   - [x] Graceful Shutdown / OS Interrupt signal handling in `main.go`.
   - [x] Idle, Read and write timeout in http.Server
   - [x] Validation of structs in [UserService](./go-chi-server/app/user/service.go) using [validator](https://github.com/go-playground/validator)

3. [https-serving](./https-serving)
   - [x]  Base [go-chi-server](./go-chi-server/)
   - [x]  HTTPS serving
   - [x]  Automatic certificate [reloading](https://opensource.com/article/22/9/dynamically-update-tls-certificates-golang-server-no-downtime) on certificate changes (e.g. renewal)

# Possible Improvements

- [ ] Config loading
- [ ] Input parameter validation/sanitization
- [ ] Parameter Object
- [ ] Context parameter
- [ ] JSON handling
- [ ] Route Versioning
- [ ] Code Quality / Static Analysis (e.g. Sonarqube, Codeclimate etc.)
- [ ] Sentry
- [ ] Feature Toggling
- [ ] UML diagram
- [ ] Exposing metrics to Prometheus using go's runtime/metrics package.
- [ ] Timing out HTTP connections
- [ ] Forcing HTTP 2.0 / 1.2
- [ ] Verify that http.ListenAndServe fields each new request on a separate goroutine by sending multiple requests and collecting runtime/metrics.
- [ ] Is it possible to print the goroutine ID? For simplifying above.
- [ ] Websocket server and client. Gorilla websocket library.
- [ ] DDoS and how to avoid it.
- [ ] Caching (local, redis etc)
