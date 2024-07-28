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
	retryMax     = 2 // Initial request + 2 reties, so total 3 requests
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
	// It's possible to convert a *retryablehttp.Client directly to a *http.Client.
	// This makes use of retryablehttp broadly applicable with minimal effort.
	sc := rc.StandardClient()

	resp, err := sc.Get(ts.URL)
	if err != nil {
		log.Fatal("error in response: ", err.Error())
	}
	defer resp.Body.Close()

	log.Printf("response code: %d\n", resp.StatusCode)
}
