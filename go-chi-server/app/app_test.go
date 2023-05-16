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

	BeforeEach(func() {
		// Create app (router)
		app := app.GetOrCreate().WithLogger(zerolog.Nop()).SetupMiddlewares().SetupCORS().SetupNotFoundHandler()
		// Create server to test the app
		ts = httptest.NewServer(app.Router)
	})

	AfterEach(func() {
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
})
