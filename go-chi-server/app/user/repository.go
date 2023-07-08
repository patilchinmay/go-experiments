package user

import (
	"context"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

var usrrepo *UserRepository

func NewUserRepository(db *gorm.DB, automigrate bool) *UserRepository {
	if usrrepo == nil {
		usrrepo = &UserRepository{
			db: db,
		}

		if automigrate {
			// automigrate the user table
			usrrepo.db.AutoMigrate(&User{})
		}
	}
	return usrrepo
}

// DiscardUserRepository will remove the reference to usrrepo so that it can be garbage collected. In other words, it deletes the singleton instance of *UserRepository.
func DiscardUserRepository() {
	if usrrepo != nil {
		usrrepo = nil
	}
}

func (ur *UserRepository) Get(ctx context.Context, id uint) (User, error) {
	var user User

	result := ur.db.First(&user, id)
	// result := ur.db.Debug().Omit("Age").First(&user, id) // Example of printing the query and ignoring a field

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

func (ur *UserRepository) Delete(ctx context.Context, id uint) error {
	result := ur.db.Delete(&User{}, id) // this is soft delete

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (ur *UserRepository) Update(ctx context.Context, id uint, input User) error {
	var user User
	result := ur.db.First(&user, id)

	if result.Error != nil {
		return result.Error
	}

	// Updates supports updating with struct or map[string]interface{}, when updating with struct it will only update non-zero fields by default
	// https://gorm.io/docs/update.html#Updates-multiple-columns
	// Only fields passed in the input will be updated. Rest will be left untouched.
	result = ur.db.Model(&user).Updates(input)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
