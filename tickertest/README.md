# Tickertest

Purpose:

1. Test `time.NewTicker` with fake time
2. Test `context.WithTimeout` with fake time

Write a program that will send a sample value `"OK"` on a `channel` at every `tick` of the `ticker`. Default `tickerInterval  = 1 * time.Second` and `timeoutInterval = 5 * time.Second`

1. Handle reading from a `channel` in `for-select`.
2. Handle `context timeout`.
3. Handle draining (and proper closure) of `channel` after `context timeout`.
4. Dependency injection with clock usage to allow mocking in tests.
5. Demonstrate `goroutine` ownership best practice.
   1. Owner of the `goroutine` writes to it and closes it.
   2. Consumer of the `goroutine` only reads it.


## Run

Notice that the two outputs of the same program differ (`draining resultCh`).
This is because the the `resultCh` may or may not have the value in it depending on the flow of the program and the moment ticker hit. More details in the code comments.

Even though the channel may have 5th value, however, main will always process only 4 values due to the nature of the program. More cases in tests.

```bash
❯ go run main.go
Tick at 2024-07-13 19:26:31.358836 +0900 JST m=+1.001266626
received: OK
Tick at 2024-07-13 19:26:32.358852 +0900 JST m=+2.001263042
received: OK
Tick at 2024-07-13 19:26:33.358627 +0900 JST m=+3.001020501
received: OK
Tick at 2024-07-13 19:26:34.358914 +0900 JST m=+4.001291501
received: OK
Tick at 2024-07-13 19:26:35.358931 +0900 JST m=+5.001293042
[function] Context expired while trying to send, exiting...
[main] context expired, waiting for resultCh to close...
[main] resultCh closed.
[main] exiting...
```

```bash
❯ go run main.go
Tick at 2024-07-13 19:28:05.250915 +0900 JST m=+1.001285917
received: OK
Tick at 2024-07-13 19:28:06.250913 +0900 JST m=+2.001285084
received: OK
Tick at 2024-07-13 19:28:07.25094 +0900 JST m=+3.001313001
received: OK
Tick at 2024-07-13 19:28:08.250936 +0900 JST m=+4.001309876
received: OK
Tick at 2024-07-13 19:28:09.251004 +0900 JST m=+5.001377917
[main] context expired, waiting for resultCh to close...
[function] context expired, exiting...
[main] STILL draining resultCh.
	Received: OK
[main] resultCh closed.
[main] exiting...
```

## Test

Since we are using mock clock for ticker and context timeout, the test will actually finish faster. We don't have to wait for the time to pass as per real clock.

```bash
❯ go test -count=1 -v ./...
=== RUN   TestExecuteOnTick
=== RUN   TestExecuteOnTick/Ticker_1s,_Timeout_5s
Tick at 1970-01-01 09:00:01 +0900 JST
Tick at 1970-01-01 09:00:02 +0900 JST
Tick at 1970-01-01 09:00:03 +0900 JST
Tick at 1970-01-01 09:00:04 +0900 JST
Tick at 1970-01-01 09:00:05 +0900 JST
[function] Context expired while trying to send, exiting...
    main_test.go:85: resultCh closed
=== RUN   TestExecuteOnTick/Ticker_2s,_Timeout_5s
Tick at 1970-01-01 09:00:02 +0900 JST
Tick at 1970-01-01 09:00:04 +0900 JST
[function] context expired, exiting...
    main_test.go:85: resultCh closed
=== RUN   TestExecuteOnTick/Ticker_2s,_Timeout_6s
Tick at 1970-01-01 09:00:02 +0900 JST
Tick at 1970-01-01 09:00:04 +0900 JST
Tick at 1970-01-01 09:00:06 +0900 JST
[function] Context expired while trying to send, exiting...
    main_test.go:85: resultCh closed
=== RUN   TestExecuteOnTick/Ticker_3s,_Timeout_10s
Tick at 1970-01-01 09:00:03 +0900 JST
Tick at 1970-01-01 09:00:06 +0900 JST
Tick at 1970-01-01 09:00:09 +0900 JST
[function] context expired, exiting...
    main_test.go:85: resultCh closed
=== RUN   TestExecuteOnTick/Ticker_5s,_Timeout_10s
Tick at 1970-01-01 09:00:05 +0900 JST
Tick at 1970-01-01 09:00:10 +0900 JST
[function] Context expired while trying to send, exiting...
    main_test.go:85: resultCh closed
--- PASS: TestExecuteOnTick (0.04s)
    --- PASS: TestExecuteOnTick/Ticker_1s,_Timeout_5s (0.01s)
    --- PASS: TestExecuteOnTick/Ticker_2s,_Timeout_5s (0.01s)
    --- PASS: TestExecuteOnTick/Ticker_2s,_Timeout_6s (0.01s)
    --- PASS: TestExecuteOnTick/Ticker_3s,_Timeout_10s (0.01s)
    --- PASS: TestExecuteOnTick/Ticker_5s,_Timeout_10s (0.01s)
PASS
ok  	github.com/patilchinmay/go-experiments/tickertest	0.299s
```