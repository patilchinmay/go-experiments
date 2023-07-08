package user_test

import (
	"context"
	"fmt"
	"math/rand"

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
		usrrepo = user.NewUserRepository(gdb.Debug(), automigrateUser)
	})

	AfterEach(func() {
		user.DiscardUserRepository()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("Get User", func() {
		It("should find a single user", func() {

			user := &user.User{
				ID:        uint(rand.Uint32()),
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

		It("should NOT find the single user", func() {

			var usrid uint = uint(rand.Uint32())
			_, err := usrrepo.Get(context.Background(), usrid)

			Expect(err).Should(HaveOccurred())
		})
	})

	Context("Add User", func() {
		It("should add a single user", func() {

			user := &user.User{
				ID:        uint(rand.Uint32()),
				FirstName: "test_firstname",
				LastName:  "test_lastname",
				Age:       30,
				Email:     "test@test.com",
			}

			dbUserID, err := usrrepo.Add(context.Background(), *user)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(dbUserID).Should(Equal(user.ID))
		})
	})

	Context("Delete User", func() {
		It("should delete a single user", func() {

			user := &user.User{
				ID:        uint(rand.Uint32()),
				FirstName: "test_firstname",
				LastName:  "test_lastname",
				Age:       30,
				Email:     "test@test.com",
			}

			// add mock data
			gdb.Create(user)

			// Should delete the user
			err := usrrepo.Delete(context.Background(), user.ID)
			Expect(err).ShouldNot(HaveOccurred())

			// Should not get deleted user
			_, err = usrrepo.Get(context.Background(), user.ID)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("Update User", func() {
		It("should update a single field, and leave others untouched", func() {

			usr := &user.User{
				ID:        uint(rand.Uint32()),
				FirstName: "test_firstname",
				LastName:  "test_lastname",
				Age:       30,
				Email:     "test@test.com",
			}

			toBeUpdated := &user.User{
				LastName: "updated_lastname",
			}

			// Add User
			dbUserID, err := usrrepo.Add(context.Background(), *usr)

			// Verify that the user is created
			Expect(err).ShouldNot(HaveOccurred())
			Expect(dbUserID).Should(Equal(usr.ID))

			// Update the user
			err = usrrepo.Update(context.Background(), usr.ID, *toBeUpdated)
			Expect(err).ShouldNot(HaveOccurred())

			// Get the updated user
			updatedUser, err := usrrepo.Get(context.Background(), usr.ID)

			fmt.Printf("CreatedAt: %v\n", updatedUser.CreatedAt)
			fmt.Printf("UpdatedAt: %v\n", updatedUser.UpdatedAt)

			// Verify that fields except lastname have NOT changed
			Expect(err).ShouldNot(HaveOccurred())
			Expect(updatedUser.ID).Should(Equal(usr.ID))
			Expect(updatedUser.FirstName).Should(Equal(usr.FirstName))
			Expect(updatedUser.Age).Should(Equal(usr.Age))
			Expect(updatedUser.Email).Should(Equal(usr.Email))

			// Verify that the lastname has changed
			Expect(updatedUser.LastName).Should(Equal(toBeUpdated.LastName))

			// Verify that the updated_at has changed
			Expect(updatedUser.UpdatedAt).ShouldNot(Equal(updatedUser.CreatedAt))
		})

		It("should update all fields", func() {

			usr := &user.User{
				ID:        uint(rand.Uint32()),
				FirstName: "test_firstname",
				LastName:  "test_lastname",
				Age:       30,
				Email:     "test@test.com",
			}

			toBeUpdated := &user.User{
				FirstName: "fn",
				LastName:  "ln",
				Age:       25,
				Email:     "fnln@test.com",
			}

			// Add User
			dbUserID, err := usrrepo.Add(context.Background(), *usr)

			// Verify that the user is created
			Expect(err).ShouldNot(HaveOccurred())
			Expect(dbUserID).Should(Equal(usr.ID))

			// Update the user
			err = usrrepo.Update(context.Background(), usr.ID, *toBeUpdated)
			Expect(err).ShouldNot(HaveOccurred())

			// Get the updated user
			updatedUser, err := usrrepo.Get(context.Background(), usr.ID)

			fmt.Printf("CreatedAt: %v\n", updatedUser.CreatedAt)
			fmt.Printf("UpdatedAt: %v\n", updatedUser.UpdatedAt)

			// Verify that ID is NOT changed
			Expect(err).ShouldNot(HaveOccurred())
			Expect(updatedUser.ID).Should(Equal(usr.ID))

			// Verify that other fields have changed
			Expect(updatedUser.FirstName).Should(Equal(toBeUpdated.FirstName))
			Expect(updatedUser.Age).Should(Equal(toBeUpdated.Age))
			Expect(updatedUser.Email).Should(Equal(toBeUpdated.Email))
			Expect(updatedUser.LastName).Should(Equal(toBeUpdated.LastName))

			// Verify that the updated_at has changed
			Expect(updatedUser.UpdatedAt).ShouldNot(Equal(updatedUser.CreatedAt))
		})
	})
})
