package testhelpers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	. "github.com/onsi/gomega"
)

// Ref: https://bismobaruno.medium.com/unit-test-http-request-in-golang-a96d146406e6
type HttpOptions struct {
	Ctx     context.Context
	Url     string
	TO      *time.Duration
	Headers map[string]string
	Queries map[string]string
	Forms   map[string]string
	Data    []byte
	Method  string
}

func DoRequest(opt *HttpOptions) (*http.Response, string) {
	if len(opt.Forms) > 0 {
		form := url.Values{}
		for key, value := range opt.Forms {
			form.Add(key, value)
		}
		opt.Data = []byte(form.Encode())
	}

	if opt.TO != nil {
		timeout := *opt.TO
		ctx, cancel := context.WithTimeout(opt.Ctx, timeout*time.Second)
		defer cancel()

		opt.Ctx = ctx
	}

	body := bytes.NewBuffer(opt.Data)
	defer body.Reset()

	req, err := http.NewRequestWithContext(opt.Ctx, opt.Method, opt.Url, body)
	if err != nil {
		Expect(err).ShouldNot(HaveOccurred())
		return nil, ""
	}
	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}

	queryValues := req.URL.Query()
	for key, val := range opt.Queries {
		queryValues.Set(key, val)
	}
	req.URL.RawQuery = queryValues.Encode()

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		Expect(err).ShouldNot(HaveOccurred())
		return nil, ""
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		Expect(err).ShouldNot(HaveOccurred())
		return nil, ""
	}

	return resp, string(respBody)
}
