# go-retryablehttp

- [go-retryablehttp](#go-retryablehttp)
  - [Description](#description)
  - [What are we building?](#what-are-we-building)
- [Run](#run)
  - [Explanation](#explanation)
  - [Reference](#reference)

## Description

As per [godoc](https://pkg.go.dev/github.com/hashicorp/go-retryablehttp):

The retryablehttp package provides a familiar HTTP client interface with automatic retries and exponential backoff. It is a thin wrapper over the standard net/http client library and exposes nearly the same public API. This makes retryablehttp very easy to drop into existing programs.

retryablehttp performs automatic retries under certain conditions. Mainly, if an error is returned by the client (connection errors, etc.), or if a 500-range response code is received (except 501), then a retry is invoked after a wait period. Otherwise, the response is returned and left to the caller to interpret.

The main difference from net/http is that requests which take a request body (POST/PUT et. al) can have the body provided in a number of ways (some more or less efficient) that allow "rewinding" the request body if the initial request fails so that the full request can be attempted again.

## What are we building?

Write a program that uses go-retryablehttp to call an api. Verify that the retries work when the api returns 5xx errors.

# Run

```bash
go run .

2024/07/28 21:34:15 [DEBUG] GET http://127.0.0.1:59832
2024/07/28 21:34:15 request no:  0
2024/07/28 21:34:15 [DEBUG] GET http://127.0.0.1:59832 (status: 500): retrying in 10ms (2 left)
2024/07/28 21:34:15 request no:  1
2024/07/28 21:34:15 [DEBUG] GET http://127.0.0.1:59832 (status: 500): retrying in 20ms (1 left)
2024/07/28 21:34:15 request no:  2
2024/07/28 21:34:15 response code: 200
```

## Explanation

The test server counts how many times it has been evoked. When the counter is less than the configured RetryMax, it responds with a 500 error response. After that, it responds with 200 OK response.

It can be seen from the above logs that the client automatically performed the 2 retries on receiving 500 error response.

It performed a total of 3 requests including the initial request + 2 retries.

## Reference

- https://github.com/hashicorp/go-retryablehttp
- http://godoc.org/github.com/hashicorp/go-retryablehttp
