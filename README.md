# go-experiments

This repository is a [go workspace](https://go.dev/doc/tutorial/workspaces) for ease of running the examples.

# How to read/understand this?

Basic familiarity of golang and backend development is assumed.

The repositories are numbered. Each one has a README.md with description, relevant details to setup and run.

Most likely each repository will build on top of the previous one, unless mentioned otherwise.

# Table of Contents

1. [hello](./hello/)
   - [x] Verify that setup works with a hello world application.
2. [go-chi-server](./go-chi-server/)
   - [x] Basic Go-chi Server (`BCS`)
   - [x] Separation of App and Server
   - [x] HTTP Request Logging (`httplog`)
   - [x] Configurable app logging (`zerolog`)
   - [x] Loads environment variables from `.env` file (`godotenv`)
   - [x] Graceful Shutdown / OS Interrupt signal handling
   - [ ] `main` handles OS interrupts
   - [ ] Builder pattern for server creation with different methods such as setTimeout
   - [ ] Tests and how to run them (`ginkgo, gomega`)
   - [ ] Explanantory comments and godoc

3.

# Possible Improvements

- [ ] Config loading
- [ ] Input parameter validation/sanitization
- [ ] Context parameter
- [ ] Middlewares
- [ ] Passing request id end-to-end in the form of `X-Correlation-ID` header
- [ ] Injectable config for logging env
- [ ] Database/ORM/Repository Layer
- [ ] JSON handling
- [ ] Route Versioning
- [ ] Tests
- [ ] Load Testing
- [ ] UML diagram
