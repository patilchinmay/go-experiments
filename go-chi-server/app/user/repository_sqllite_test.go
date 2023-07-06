package user_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/patilchinmay/go-experiments/go-chi-server/app/user"
)

var _ = Describe("UserRepository with SQLite", func() {
	var usrrepo *user.UserRepository
	var gdb *gorm.DB
	var err error

	BeforeEach(func() {
		// open gorm db
		// When the name of the database file handed to sqlite3_open() or to ATTACH is an empty string, then a new temporary file is created to hold the database.
		// https://www.sqlite.org/inmemorydb.html
		gdb, err = gorm.Open(sqlite.Open(""), &gorm.Config{})
		Expect(err).ShouldNot(HaveOccurred())

		automigrateUser := true
		usrrepo = user.NewUserRepository(gdb, automigrateUser)
	})

	AfterEach(func() {
		user.DiscardUserRepository()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("Get User", func() {
		It("should find a single user", func() {

			user := &user.User{
				ID:        100,
				FirstName: "test_firstname",
				LastName:  "test_lastname",
				Age:       30,
				Email:     "test@test.com",
			}

			// add mock data
			gdb.Create(user)

			dbUser, err := usrrepo.Get(context.Background(), user.ID)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(dbUser.ID).Should(Equal(user.ID))
		})
	})
})
