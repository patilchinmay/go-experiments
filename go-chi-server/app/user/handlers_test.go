package user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patilchinmay/go-experiments/go-chi-server/app"
	"github.com/patilchinmay/go-experiments/go-chi-server/app/user"
	"github.com/patilchinmay/go-experiments/go-chi-server/utils/testhelpers"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _ = Describe("User Handlers", Serial, func() {
	var ts *httptest.Server
	var App *app.App
	var gdb *gorm.DB
	var path string = "/user"

	BeforeEach(func() {
		// open gorm db
		// When the name of the database file handed to sqlite3_open() or to ATTACH is an empty string, then a new temporary file is created to hold the database.
		// https://www.sqlite.org/inmemorydb.html
		var err error
		gdb, err = gorm.Open(sqlite.Open(""), &gorm.Config{})
		Expect(err).ShouldNot(HaveOccurred())

		// logger
		logger := zerolog.Nop()
		// Create app with routes handlers (uses builder pattern)
		App = app.GetOrCreate().SetupDB(gdb).WithLogger(logger).SetupCORS().SetupMiddlewares().SetupNotFoundHandler()

		// Create and setup user subrouter
		user.SetupSubrouter(gdb, logger)

		// Initialize and register subrouters
		App.MountSubrouters()

		// Create server to test the app
		ts = httptest.NewServer(App.Router)
	})

	AfterEach(func() {
		ts.Close()
		app.Discard()
		gdb = nil
	})

	// Path /user
	Context("Get User Handler", func() {

		// GET /user
		When("User ID is invalid", func() {
			It("Should return http 400 error", func() {
				to := time.Duration(10)
				opt := &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path + "/abc",
					TO:     &to,
					Method: http.MethodGet,
				}

				res, _ := testhelpers.DoRequest(opt)

				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				// Expect(bodystring).To(ContainSubstring(`"Ping":"Pong"`))
				// Expect(res).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
			})
		})
	})
})
