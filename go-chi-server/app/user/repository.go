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

func (ur *UserRepository) Get(ctx context.Context, id string) (User, error) {
	var user User

	result := ur.db.Debug().Omit("Age").First(&user, id) // Example of printing the query and ignoring a field

	if result.Error != nil {
		return user, result.Error
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

func (ur *UserRepository) Delete(ctx context.Context, id string) error {
	result := ur.db.Delete(&User{}, id) // this is soft delete

	if result.Error != nil {
		return result.Error
	}

	return nil
}
