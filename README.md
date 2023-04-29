# go-experiments

This repository is a [go workspace](https://go.dev/doc/tutorial/workspaces) for ease of running the examples.

# How to read/understand this?

The repositories are numbered. Each one has a README.md with description, relevant details to setup and run. Most likely each repository will build on top of the previous one, unless mentioned otherwise.

# ToC

1. [hello](./hello/)
   - [x] Verify that setup works with a hello world application.
2. [go-chi-server](./go-chi-server/)
   - [x] Basic Go-chi Server (BCS)
   - [x] Separation of App and Server
   - [x] HTTP Request Logging (httplog)
   - [x] Configurable app logging (zerolog)
   - [x] Loads environment variables from `.env` file (godotenv)
3.

# Possible Improvements

- [ ] Explanantory comments and godoc
- [ ] Interrupt signal handling
- [ ] Config Loading
- [ ] Input parameter validation/sanitization
- [ ] Context
- [ ] Passing request id end-to-end in the form of `X-Correlation-ID` header
- [ ] Injectable config for logging env
- [ ] Database/ORM/Repository Layer
- [ ] JSON handling
- [ ] Tests
- [ ] Load Testing
- [ ] UML diagram
