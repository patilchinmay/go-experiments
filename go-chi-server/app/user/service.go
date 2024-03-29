package user

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	cnp "github.com/patilchinmay/go-experiments/cloudnativepatterns"
	"github.com/patilchinmay/go-experiments/go-chi-server/utils/logger"
	v "github.com/patilchinmay/go-experiments/go-chi-server/utils/validator"
)

type UserService struct {
	usrrepo UserRepository
	cnp     cnp.CloudNativePatterns
}

var usrsvc *UserService

func NewUserService(usrrepo UserRepository, cnp cnp.CloudNativePatterns) *UserService {
	if usrsvc == nil {
		usrsvc = &UserService{
			usrrepo: usrrepo,
			cnp:     cnp,
		}
	}
	return usrsvc
}

// DiscardUserService will remove the reference to usrsvc so that it can be garbage collected. In other words, it deletes the singleton instance of *UserService.
func DiscardUserService() {
	if usrsvc != nil {
		usrsvc = nil
	}
}

const maxRetries = 3
const retryInterval = 2 * time.Second

func (u *UserService) Get(ctx context.Context, id uint) (User, error) {
	logger := logger.Logger.With().Str("requestID", middleware.GetReqID(ctx)).Logger()
	logger.Debug().Msg("User Service : Get")

	var (
		user User
		err  error
	)

	// Retry-able function
	var userRepoGet cnp.CloudNativeFunction = func(ctx context.Context) error {
		// Call the repository layer
		user, err = u.usrrepo.Get(ctx, id)
		if err != nil {
			return err
		}

		return nil
	}

	r := u.cnp.Retry(userRepoGet, maxRetries, retryInterval)

	err = r(ctx)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *UserService) Add(ctx context.Context, user User) (uint, error) {
	//  Setup logger
	logger := logger.Logger.With().Str("requestID", middleware.GetReqID(ctx)).Logger()
	logger.Debug().Msg("User Service : Add")

	// Validate the input
	err := v.Validator.Struct(user)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Error().Err(err).Msg("Failed validation")
			return 0, err
		}

		for _, err := range err.(validator.ValidationErrors) {
			logger.Error().Err(err).Msg(err.Error())
		}

		// from here you can create your own error messages in whatever language you wish
		return 0, err
	}

	var id uint

	// Retry-able function
	var userRepoAdd cnp.CloudNativeFunction = func(ctx context.Context) error {
		// Call the repository layer
		id, err = u.usrrepo.Add(ctx, user)
		if err != nil {
			return err
		}

		return nil
	}

	r := u.cnp.Retry(userRepoAdd, maxRetries, retryInterval)

	err = r(ctx)

	if err != nil {
		return 0, err
	}

	// return the response
	return id, nil
}

func (u *UserService) Delete(ctx context.Context, id uint) error {
	logger := logger.Logger.With().Str("requestID", middleware.GetReqID(ctx)).Logger()
	logger.Debug().Msg("User Service : Delete")

	var err error

	// Retry-able function
	var userRepoDelete cnp.CloudNativeFunction = func(ctx context.Context) error {
		// Call the repository layer
		err = u.usrrepo.Delete(ctx, id)
		if err != nil {
			return err
		}

		return nil
	}

	r := u.cnp.Retry(userRepoDelete, maxRetries, retryInterval)

	err = r(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) Update(ctx context.Context, id uint, input UpdateUserInput) error {
	logger := logger.Logger.With().Str("requestID", middleware.GetReqID(ctx)).Logger()
	logger.Debug().Msg("User Service : Update")

	// Validate the input
	err := v.Validator.Struct(input)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Error().Err(err).Msg("Failed validation")
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {
			logger.Error().Err(err).Msg(err.Error())
		}

		// from here you can create your own error messages in whatever language you wish
		return err
	}

	var user = User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Age:       input.Age,
		Email:     input.Email,
	}

	// Retry-able function
	var userRepoUpdate cnp.CloudNativeFunction = func(ctx context.Context) error {
		// Call the repository layer
		err = u.usrrepo.Update(ctx, id, user)
		if err != nil {
			return err
		}

		return nil
	}

	r := u.cnp.Retry(userRepoUpdate, maxRetries, retryInterval)

	err = r(ctx)

	if err != nil {
		return err
	}

	return nil
}
