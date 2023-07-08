package user_test

import (
	"context"
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
				Add(context.Background(), usr).Return(usr.ID, nil)

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

})
