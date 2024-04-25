# Testcontainers

## Description

This is a demo of using [testcontainers](https://testcontainers.com/).

We create a package `testcontainers-demo`.

It uses the redis client and has 2 simple functions viz. Set and Get

To run this program, we would need a live instance/container of redis/valkey as provided in the `docker-compose.yaml` file.

However, for testing purposes, creating a container manually is cumbersome and not recommended.

[testcontainers](https://testcontainers.com/) helps spin up a valkey container for the testing purpose via code.

## Test

```bash
❯ go test ./... -v

=== RUN   TestWithRedis
2024/04/25 17:07:02 github.com/testcontainers/testcontainers-go - Connected to docker:
  Server Version: 25.0.3
  API Version: 1.44
  Operating System: Docker Desktop
  Total Memory: 7841 MB
  Resolved Docker Host: unix:///var/run/docker.sock
  Resolved Docker Socket Path: /var/run/docker.sock
  Test SessionID: ddea0335b69d2950f647140154032966ca138e52295579b4d11235083e0312f3
  Test ProcessID: af98ebe4-90cb-48f1-b658-33d411df5cee
2024/04/25 17:07:02 🐳 Creating container for image testcontainers/ryuk:0.7.0
2024/04/25 17:07:02 ✅ Container created: 3c1875fd55f8
2024/04/25 17:07:02 🐳 Starting container: 3c1875fd55f8
2024/04/25 17:07:02 ✅ Container started: 3c1875fd55f8
2024/04/25 17:07:02 🚧 Waiting for container id 3c1875fd55f8 image: testcontainers/ryuk:0.7.0. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms}
2024/04/25 17:07:03 🔔 Container is ready: 3c1875fd55f8
2024/04/25 17:07:03 🐳 Creating container for image valkey/valkey:7.2.5
2024/04/25 17:07:03 ✅ Container created: 95504a5d8691
2024/04/25 17:07:03 🐳 Starting container: 95504a5d8691
2024/04/25 17:07:03 ✅ Container started: 95504a5d8691
2024/04/25 17:07:03 🚧 Waiting for container id 95504a5d8691 image: valkey/valkey:7.2.5. Waiting for: &{timeout:<nil> Log:Ready to accept connections IsRegexp:false Occurrence:1 PollInterval:100ms}
2024/04/25 17:07:03 🔔 Container is ready: 95504a5d8691
=== RUN   TestWithRedis/Set_foo_bar
=== RUN   TestWithRedis/Get_foo
foo bar
2024/04/25 17:07:03 🐳 Terminating container: 95504a5d8691
2024/04/25 17:07:03 🚫 Container terminated: 95504a5d8691
--- PASS: TestWithRedis (0.75s)
    --- PASS: TestWithRedis/Set_foo_bar (0.01s)
    --- PASS: TestWithRedis/Get_foo (0.00s)
PASS
ok  	github.com/patilchinmay/go-experiments/testcontainersdemo	1.258s
```

#### Reference

- https://www.youtube.com/watch?v=sNg0bnMF_qY
- https://github.com/dreamsofcode-io/testcontainers/blob/main/pubsub/pubsub_test.go
- https://golang.testcontainers.org/quickstart/