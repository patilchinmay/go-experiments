package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

var (
	retryCounter = 0
	retryMax     = 9 // Initial request + 9 reties, so total 10 requests
)

func main() {
	// Set up test server
	// It counts the number of retries and responds accordingly
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("request no: ", retryCounter)

		if retryCounter < retryMax {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		retryCounter++
	}))
	defer ts.Close()

	rc := retryablehttp.NewClient()

	rc.RetryMax = retryMax
	rc.RetryWaitMin = 10 * time.Millisecond
	rc.RetryWaitMax = 50 * time.Millisecond

	// Optional
	// Test if retries work after substituting the HTTPClient
	// This maybe useful if want to use a different HTTPClient
	// Such as one with token auto-refresh capability
	// e.g. https://github.com/patilchinmay/go-experiments/tree/master/http-client-autorefreshtoken
	// rc.HTTPClient.Transport = http.DefaultTransport
	// Finding: The retries work even after substituting the HTTPClient
	rc.HTTPClient = http.DefaultClient

	// Optional
	// It's possible to convert a *retryablehttp.Client directly to a *http.Client.
	// This makes use of retryablehttp broadly applicable with minimal effort.
	// sc := rc.StandardClient()

	resp, err := rc.Get(ts.URL)
	if err != nil {
		log.Fatal("error in response: ", err.Error())
	}
	defer resp.Body.Close()

	log.Printf("response code: %d\n", resp.StatusCode)
}
