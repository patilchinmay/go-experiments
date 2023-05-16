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

	BeforeEach(func() {
		// Create app
		app := app.GetOrCreate().WithLogger(zerolog.Nop()).SetupMiddlewares().SetupCORS().SetupNotFoundHandler()

		// We do not need to instantiate ping as tests will implicitly run the init function of ping package which will instantiate itself

		// Initialize and register subrouters
		app.MountSubrouters()

		// Create server to test the app
		ts = httptest.NewServer(app.Router)
	})

	AfterEach(func() {
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
			Expect(bodystring).To(Equal(`{"Ping":"Pong"}`))
		})
	})
})
