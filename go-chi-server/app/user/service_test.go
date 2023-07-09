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

		It("should return an error if failed to delete", func() {
			userID := uint(rand.Uint32())

			// Define mock expectation
			dummyError := errors.New("dummy error")
			usrrepomock.
				EXPECT().
				Delete(context.Background(), userID).
				Return(dummyError).
				Times(1)

			err := usrsvc.Delete(context.Background(), userID)

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("dummy error"))
		})
	})

})
