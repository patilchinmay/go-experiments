package user

import (
	"context"

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

func (ur *UserRepository) Add(ctx context.Context, user User) (uint, error) {
	result := ur.db.Create(&user)

	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}
