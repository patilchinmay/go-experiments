package user_test

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/patilchinmay/go-experiments/go-chi-server/app/user"
)

var _ = Describe("UserRepository with go-sqlmock", func() {
	var usrrepo user.UserRepository
	var mock sqlmock.Sqlmock

	BeforeEach(func() {
		var db *sql.DB
		var err error

		db, mock, err = sqlmock.New() // mock sql.DB
		Expect(err).ShouldNot(HaveOccurred())

		gdb, err := gorm.Open(postgres.New(
			postgres.Config{
				Conn:       db,
				DriverName: "postgres",
			},
		), &gorm.Config{}) // open gorm db
		Expect(err).ShouldNot(HaveOccurred())

		automigrateUser := false
		usrrepo = user.NewUserRepository(gdb, automigrateUser)
	})

	AfterEach(func() {
		user.DiscardUserRepository()
		err := mock.ExpectationsWereMet() // make sure all expectations were met
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
			rows := sqlmock.
				NewRows([]string{"id", "firstname", "lastname", "age", "email"}).
				AddRow(user.ID, user.FirstName, user.LastName, user.Age, user.Email)

			const sql = `SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`

			mock.ExpectQuery(regexp.QuoteMeta(sql)).WithArgs(user.ID).WillReturnRows(rows)

			dbUser, err := usrrepo.Get(context.Background(), user.ID)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(dbUser.ID).Should(Equal(user.ID))
		})
	})
})
