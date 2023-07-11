package user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
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
		ts = nil
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

		// GET /user
		When("User ID is of non-existent user", func() {
			It("Should return http 500 error", func() {
				to := time.Duration(10)
				opt := &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path + "/12345",
					TO:     &to,
					Method: http.MethodGet,
				}

				res, _ := testhelpers.DoRequest(opt)

				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				// Expect(bodystring).To(ContainSubstring(`"Ping":"Pong"`))
				// Expect(res).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
			})
		})
	})

	Context("Add User Handler", func() {

		// POST /user
		When("Valid body is present", func() {
			It("Should add user without an error", func() {
				body := []byte(`{
					"firstname": "abc",
					"lastname": "xyz",
					"age": 29,
					"email": "abcxyz@test.com"
				}`)
				to := time.Duration(10)
				opt := &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path,
					TO:     &to,
					Method: http.MethodPost,
					Data:   body,
				}

				res, bodystring := testhelpers.DoRequest(opt)

				userid, _ := strconv.Atoi(bodystring)

				Expect(res.StatusCode).To(Equal(http.StatusCreated))
				Expect(res).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
				Expect(userid).To(BeAssignableToTypeOf(0))
			})
		})

		// POST /user
		When("Malformed body is present", func() {
			It("Should return a 400 error", func() {
				// Missing a comma after "age"
				body := []byte(`{
					"firstname": "abc",
					"lastname": "xyz",
					"age": 29
					"email": "abcxyz@test.com"
				}`)
				to := time.Duration(10)
				opt := &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path,
					TO:     &to,
					Method: http.MethodPost,
					Data:   body,
				}

				res, _ := testhelpers.DoRequest(opt)

				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
				// Expect(bodystring).To(ContainSubstring(`"Ping":"Pong"`))
				// Expect(res).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
			})
		})

		// POST /user
		When("Incorrect body param is present", func() {
			It("Should return a 500 error", func() {
				// age is incorrect. Should be 0 >= age >= 130
				body := []byte(`{
					"firstname": "abc",
					"lastname": "xyz",
					"age": 200,
					"email": "abcxyz@test.com"
				}`)
				to := time.Duration(10)
				opt := &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path,
					TO:     &to,
					Method: http.MethodPost,
					Data:   body,
				}

				res, _ := testhelpers.DoRequest(opt)

				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				// Expect(bodystring).To(ContainSubstring(`"Ping":"Pong"`))
				// Expect(res).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
			})
		})
	})

	Context("Delete User Handler", func() {

		// DELETE /user
		When("Valid user id is present", func() {
			It("Should delete user without an error", func() {
				body := []byte(`{
					"firstname": "abc",
					"lastname": "xyz",
					"age": 29,
					"email": "abcxyz@test.com"
					}`)

				to := time.Duration(10)
				opt := &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path,
					TO:     &to,
					Method: http.MethodPost,
					Data:   body,
				}

				res, bodystring := testhelpers.DoRequest(opt)

				userid, _ := strconv.Atoi(bodystring)

				Expect(res.StatusCode).To(Equal(http.StatusCreated))
				Expect(res).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
				Expect(userid).To(BeAssignableToTypeOf(0))

				opt = &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path + "/" + bodystring,
					TO:     &to,
					Method: http.MethodDelete,
				}

				res, _ = testhelpers.DoRequest(opt)

				Expect(res).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		// DELETE /user
		When("User ID is invalid", func() {
			It("Should return http 400 error", func() {
				to := time.Duration(10)
				opt := &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path + "/abc",
					TO:     &to,
					Method: http.MethodDelete,
				}

				res, _ := testhelpers.DoRequest(opt)

				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Context("Update User Handler", func() {

		// PATCH /user
		When("Valid user id is present", func() {
			It("Should update user without an error", func() {
				body := []byte(`{
					"firstname": "abc",
					"lastname": "xyz",
					"age": 29,
					"email": "abcxyz@test.com"
					}`)

				to := time.Duration(10)
				opt := &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path,
					TO:     &to,
					Method: http.MethodPost,
					Data:   body,
				}

				res, bodystring := testhelpers.DoRequest(opt)

				userid, _ := strconv.Atoi(bodystring)

				Expect(res.StatusCode).To(Equal(http.StatusCreated))
				Expect(res).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
				Expect(userid).To(BeAssignableToTypeOf(0))

				body = []byte(`{
					"firstname": "updatefn",
					"lastname": "updateln",
					"age": 35,
					"email": "update@test.com"
					}`)

				opt = &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path + "/" + bodystring,
					TO:     &to,
					Method: http.MethodPatch,
					Data:   body,
				}

				res, _ = testhelpers.DoRequest(opt)

				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		// PATCH /user
		When("User ID is invalid", func() {
			It("Should return http 400 error", func() {
				body := []byte(`{
					"firstname": "updatefn",
					"lastname": "updateln",
					"age": 35,
					"email": "update@test.com"
					}`)
				to := time.Duration(10)
				opt := &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path + "/abc",
					TO:     &to,
					Method: http.MethodPatch,
					Data:   body,
				}

				res, _ := testhelpers.DoRequest(opt)

				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		// PATCH /user
		When("User ID is of non-existent user", func() {
			It("Should return http 500 error", func() {
				body := []byte(`{
					"firstname": "updatefn",
					"lastname": "updateln",
					"age": 35,
					"email": "update@test.com"
					}`)
				to := time.Duration(10)
				opt := &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path + "/12345",
					TO:     &to,
					Method: http.MethodPatch,
					Data:   body,
				}

				res, _ := testhelpers.DoRequest(opt)

				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})

		// PATCH /user
		When("Malformed body is present", func() {
			It("Should return a 400 error", func() {
				body := []byte(`{
					"firstname": "abc",
					"lastname": "xyz",
					"age": 29,
					"email": "abcxyz@test.com"
					}`)

				to := time.Duration(10)
				opt := &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path,
					TO:     &to,
					Method: http.MethodPost,
					Data:   body,
				}

				res, bodystring := testhelpers.DoRequest(opt)

				userid, _ := strconv.Atoi(bodystring)

				Expect(res.StatusCode).To(Equal(http.StatusCreated))
				Expect(res).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
				Expect(userid).To(BeAssignableToTypeOf(0))

				// Missing a comma after "age"
				body = []byte(`{
					"firstname": "updatefn",
					"lastname": "updateln",
					"age": 35
					"email": "update@test.com"
					}`)
				to = time.Duration(10)
				opt = &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path + "/" + bodystring,
					TO:     &to,
					Method: http.MethodPatch,
					Data:   body,
				}

				res, _ = testhelpers.DoRequest(opt)

				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		// PATCH /user
		When("Incorrect body param is present", func() {
			It("Should return a 500 error", func() {
				body := []byte(`{
					"firstname": "abc",
					"lastname": "xyz",
					"age": 29,
					"email": "abcxyz@test.com"
					}`)

				to := time.Duration(10)
				opt := &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path,
					TO:     &to,
					Method: http.MethodPost,
					Data:   body,
				}

				res, bodystring := testhelpers.DoRequest(opt)

				userid, _ := strconv.Atoi(bodystring)

				Expect(res.StatusCode).To(Equal(http.StatusCreated))
				Expect(res).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
				Expect(userid).To(BeAssignableToTypeOf(0))

				// age is incorrect. Should be 0 >= age >= 130
				body = []byte(`{
					"firstname": "updatefn",
					"lastname": "updateln",
					"age": 200,
					"email": "update@test.com"
					}`)
				to = time.Duration(10)
				opt = &testhelpers.HttpOptions{
					Ctx:    context.Background(),
					Url:    ts.URL + path + "/" + bodystring,
					TO:     &to,
					Method: http.MethodPatch,
					Data:   body,
				}

				res, _ = testhelpers.DoRequest(opt)

				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})

	})
})
