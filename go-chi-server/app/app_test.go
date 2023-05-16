package app_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"

	"github.com/patilchinmay/go-experiments/go-chi-server/app"
	"github.com/patilchinmay/go-experiments/go-chi-server/utils/testhelpers"
)

var _ = Describe("App", func() {
	var ts = &httptest.Server{}
	var App *app.App

	BeforeEach(func() {
		// Create app with routes handlers (uses builder pattern)
		App = app.GetOrCreate().WithLogger(zerolog.Nop()).SetupCORS().SetupMiddlewares().SetupNotFoundHandler()

		// Create server to test the app
		ts = httptest.NewServer(App.Router)
	})

	AfterEach(func() {
		// Since *App is singleton, the tests would not be able to create a new instance of *App for each It().
		// Thus, we need to delete the singleton *App, so that each test It() can generate a fresh instance of *App for testing.
		defer app.Discard()
		defer ts.Close()
	})

	// GET /health healthCheck
	Context("Healthcheck", func() {
		It("should return http 200 success", func() {
			to := time.Duration(10)
			opt := &testhelpers.HttpOptions{
				Ctx:    context.Background(),
				Url:    ts.URL + "/health",
				TO:     &to,
				Method: http.MethodGet,
			}

			res, _ := testhelpers.DoRequest(opt)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		})
	})

	// GET /404
	Context("NotFoundHandler ", func() {
		It("should return http 404", func() {
			to := time.Duration(10)
			opt := &testhelpers.HttpOptions{
				Ctx:    context.Background(),
				Url:    ts.URL + "/404",
				TO:     &to,
				Method: http.MethodGet,
			}

			res, _ := testhelpers.DoRequest(opt)
			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should return http 404", func() {
			to := time.Duration(10)
			opt := &testhelpers.HttpOptions{
				Ctx:    context.Background(),
				Url:    ts.URL + "/abcd",
				TO:     &to,
				Method: http.MethodGet,
			}

			res, _ := testhelpers.DoRequest(opt)
			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})

	})

})
