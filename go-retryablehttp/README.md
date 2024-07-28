# go-retryablehttp

- [go-retryablehttp](#go-retryablehttp)
  - [Description](#description)
  - [What are we building?](#what-are-we-building)
  - [Run](#run)
  - [Explanation](#explanation)
  - [Note](#note)
  - [Reference](#reference)

## Description

As per [godoc](https://pkg.go.dev/github.com/hashicorp/go-retryablehttp):

The retryablehttp package provides a familiar HTTP client interface with automatic retries and exponential backoff. It is a thin wrapper over the standard net/http client library and exposes nearly the same public API. This makes retryablehttp very easy to drop into existing programs.

retryablehttp performs automatic retries under certain conditions. Mainly, if an error is returned by the client (connection errors, etc.), or if a 500-range response code is received (except 501), then a retry is invoked after a wait period. Otherwise, the response is returned and left to the caller to interpret.

The main difference from net/http is that requests which take a request body (POST/PUT et. al) can have the body provided in a number of ways (some more or less efficient) that allow "rewinding" the request body if the initial request fails so that the full request can be attempted again.

## What are we building?

Write a program that uses go-retryablehttp to call an api. Verify that the retries work when the api returns 5xx errors.

## Run

```bash
go run .

2024/07/28 22:02:48 [DEBUG] GET http://127.0.0.1:60307
2024/07/28 22:02:48 request no:  0
2024/07/28 22:02:48 [DEBUG] GET http://127.0.0.1:60307 (status: 500): retrying in 10ms (9 left)
2024/07/28 22:02:48 request no:  1
2024/07/28 22:02:48 [DEBUG] GET http://127.0.0.1:60307 (status: 500): retrying in 20ms (8 left)
2024/07/28 22:02:48 request no:  2
2024/07/28 22:02:48 [DEBUG] GET http://127.0.0.1:60307 (status: 500): retrying in 40ms (7 left)
2024/07/28 22:02:48 request no:  3
2024/07/28 22:02:48 [DEBUG] GET http://127.0.0.1:60307 (status: 500): retrying in 50ms (6 left)
2024/07/28 22:02:48 request no:  4
2024/07/28 22:02:48 [DEBUG] GET http://127.0.0.1:60307 (status: 500): retrying in 50ms (5 left)
2024/07/28 22:02:48 request no:  5
2024/07/28 22:02:48 [DEBUG] GET http://127.0.0.1:60307 (status: 500): retrying in 50ms (4 left)
2024/07/28 22:02:48 request no:  6
2024/07/28 22:02:48 [DEBUG] GET http://127.0.0.1:60307 (status: 500): retrying in 50ms (3 left)
2024/07/28 22:02:48 request no:  7
2024/07/28 22:02:48 [DEBUG] GET http://127.0.0.1:60307 (status: 500): retrying in 50ms (2 left)
2024/07/28 22:02:48 request no:  8
2024/07/28 22:02:48 [DEBUG] GET http://127.0.0.1:60307 (status: 500): retrying in 50ms (1 left)
2024/07/28 22:02:48 request no:  9
2024/07/28 22:02:48 response code: 200
```

## Explanation

The test server counts how many times it has been evoked. When the counter is less than the configured `RetryMax`, it responds with a `500` error response. After that, it responds with `200` OK response.

It can be seen from the above logs that the client automatically performed the 9 retries on receiving `500` error response.

It performed a total of 10 requests including the initial request + 9 retries.

## Note

If we are using the `*retryablehttp.Client`, we can substitute the `retryablehttp.Client.HTTPClient` and the retries will still work.

This maybe useful if want to use a different `HTTPClient`, such as one with token auto-refresh capability.
e.g. https://github.com/patilchinmay/go-experiments/tree/master/http-client-autorefreshtoken

This does *NOT* apply if we are using the (standard) converted client using `*retryablehttp.Client.StandardClient()`.

## Reference

- https://github.com/hashicorp/go-retryablehttp
- http://godoc.org/github.com/hashicorp/go-retryablehttp
