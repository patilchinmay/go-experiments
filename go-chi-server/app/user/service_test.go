package user_test

import (
	"context"
	"errors"
	"math/rand"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patilchinmay/go-experiments/go-chi-server/app/user"
	"github.com/patilchinmay/go-experiments/go-chi-server/app/user/mocks"
	"go.uber.org/mock/gomock"
)

var _ = Describe("User Service", func() {
	var (
		ctrl        *gomock.Controller
		usrrepomock *mocks.MockUserRepository
		usrsvc      *user.UserService
	)

	BeforeEach(func() {
		// Define mock controller
		ctrl = gomock.NewController(GinkgoT())

		// Create mock UserRepo
		usrrepomock = mocks.NewMockUserRepository(ctrl)

		// Instantiate UserService
		usrsvc = user.NewUserService(usrrepomock)
	})

	AfterEach(func() {
		ctrl.Finish()
		user.DiscardUserService()
	})

	Context("Add User", func() {
		It("should add a single user without error", func() {
			usr := user.User{
				ID:        uint(rand.Uint32()),
				FirstName: "test_firstname",
				LastName:  "test_lastname",
				Age:       30,
				Email:     "test@test.com",
			}

			// Define mock expe
			usrrepomock.
				EXPECT().
				Add(context.Background(), usr).
				Return(usr.ID, nil).
				Times(1)

			expectedUserID, err := usrsvc.Add(context.Background(), usr)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(expectedUserID).Should(Equal(usr.ID))
		})

		It("should retry 3 times for adding a single user when an error occurs", func() {
			usr := user.User{
				ID:        uint(rand.Uint32()),
				FirstName: "test_firstname",
				LastName:  "test_lastname",
				Age:       30,
				Email:     "test@test.com",
			}

			// Define mock expe
			dummyError := errors.New("")
			usrrepomock.
				EXPECT().
				Add(context.Background(), usr).
				Return(uint(0), dummyError).
				Times(4)

			_, err := usrsvc.Add(context.Background(), usr)

			Expect(err).Should(HaveOccurred())
		})

		It("should return a validation error for invalid input", func() {
			usr := user.User{
				ID:        uint(rand.Uint32()),
				FirstName: "test_firstname",
				LastName:  "test_lastname",
				Age:       250, // invalid input
				Email:     "test@test.com",
			}

			expectedUserID, err := usrsvc.Add(context.Background(), usr)

			Expect(err).Should(HaveOccurred())
			Expect(expectedUserID).Should(Equal(uint(0)))
		})
	})

	Context("Delete User", func() {
		It("should delete a single user without error", func() {
			userID := uint(rand.Uint32())

			// Define mock expectation
			usrrepomock.
				EXPECT().
				Delete(context.Background(), userID).
				Return(nil).
				Times(1)

			err := usrsvc.Delete(context.Background(), userID)

			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return an error if failed to delete after 3 retries", func() {
			userID := uint(rand.Uint32())

			// Define mock expectation
			dummyError := errors.New("dummy error")
			usrrepomock.
				EXPECT().
				Delete(context.Background(), userID).
				Return(dummyError).
				Times(4)

			err := usrsvc.Delete(context.Background(), userID)

			Expect(err).Should(HaveOccurred())
		})
	})

	Context("Get User", func() {
		It("should get a single user without error", func() {
			usr := user.User{
				ID:        uint(rand.Uint32()),
				FirstName: "test_firstname",
				LastName:  "test_lastname",
				Age:       25,
				Email:     "test@test.com",
			}
			// Define mock expectation
			usrrepomock.
				EXPECT().Get(context.Background(), usr.ID).
				Return(usr, nil).
				Times(1)

			getUserResult, err := usrsvc.Get(context.Background(), usr.ID)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(getUserResult).To(Equal(usr))
		})

		It("should return an error if failed to get user", func() {
			userID := uint(rand.Uint32())

			// Define mock expectation
			dummyError := errors.New("")
			var usr user.User

			usrrepomock.
				EXPECT().
				Get(context.Background(), userID).
				Return(usr, dummyError).
				Times(4)

			getUserResult, err := usrsvc.Get(context.Background(), userID)

			Expect(err).Should(HaveOccurred())
			Expect(getUserResult).To(Equal(usr))
		})
	})

	Context("Update User", func() {
		It("should update a single user without error", func() {
			userID := uint(rand.Uint32())

			input := user.UpdateUserInput{
				FirstName: "updatedfn",
				LastName:  "updatedln",
				Age:       25,
				Email:     "update@test.com",
			}

			usr := user.User{
				FirstName: input.FirstName,
				LastName:  input.LastName,
				Age:       input.Age,
				Email:     input.Email,
			}

			// Define mock expectation
			usrrepomock.
				EXPECT().Update(context.Background(), userID, usr).
				Return(nil).
				Times(1)

			err := usrsvc.Update(context.Background(), userID, input)

			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return an error when failed to update a single user due to validation", func() {
			userID := uint(rand.Uint32())

			input := user.UpdateUserInput{
				FirstName: "updatedfn",
				LastName:  "updatedln",
				Age:       250, // invalid input
				Email:     "update@test.com",
			}

			err := usrsvc.Update(context.Background(), userID, input)

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("validation for 'Age' failed"))
		})

		It("should retry 3 times and return an error when failed to update a single user due to repo error", func() {
			userID := uint(rand.Uint32())

			input := user.UpdateUserInput{
				FirstName: "updatedfn",
				LastName:  "updatedln",
				Age:       25,
				Email:     "update@test.com",
			}

			usr := user.User{
				FirstName: input.FirstName,
				LastName:  input.LastName,
				Age:       input.Age,
				Email:     input.Email,
			}

			// Define mock expectation
			dummyError := errors.New("")

			// Define mock expectation
			usrrepomock.
				EXPECT().Update(context.Background(), userID, usr).
				Return(dummyError).
				Times(4)

			err := usrsvc.Update(context.Background(), userID, input)

			Expect(err).Should(HaveOccurred())
		})
	})
})
