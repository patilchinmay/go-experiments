package user

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

var usrrepo *UserRepository

func NewUserRepository(db *gorm.DB) *UserRepository {
	if usrrepo == nil {
		usrrepo = &UserRepository{
			db: db,
		}

		// automigrate the user table
		usrrepo.db.AutoMigrate(&User{})
	}
	return usrrepo
}

func (ur *UserRepository) Get(ctx context.Context, id uint) (User, error) {
	var user User

	ur.db.Debug().Omit("Age").First(&user, id) // Example of printing the query and ignoring a field

	if (user == User{}) {
		return User{}, errors.New("Not found")
	}

	return user, nil
}

func (ur *UserRepository) Add(ctx context.Context, user User) (uint, error) {
	result := ur.db.Create(&user)

	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}
