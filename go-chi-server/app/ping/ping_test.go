package ping_test

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

var _ = Describe("Ping", func() {
	var ts = &httptest.Server{}
	var App *app.App

	BeforeEach(func() {
		// Create app with routes handlers (uses builder pattern)
		App = app.GetOrCreate().WithLogger(zerolog.Nop()).SetupCORS().SetupMiddlewares().SetupNotFoundHandler()

		// We do not need to instantiate ping as tests will implicitly run the init function of ping package which will instantiate itself

		// Initialize and register subrouters
		App.MountSubrouters()

		// Create server to test the app
		ts = httptest.NewServer(App.Router)
	})

	AfterEach(func() {
		defer app.Discard()
		defer ts.Close()
	})

	// Path /ping
	Context("Ping", func() {

		// GET /ping
		It("should return http 200 success", func() {
			to := time.Duration(10)
			opt := &testhelpers.HttpOptions{
				Ctx:    context.Background(),
				Url:    ts.URL + "/ping",
				TO:     &to,
				Method: http.MethodGet,
			}

			res, bodystring := testhelpers.DoRequest(opt)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(bodystring).To(ContainSubstring(`"Ping":"Pong"`))
			Expect(res).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
		})
	})
})
