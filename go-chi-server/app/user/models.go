package user

import (
	"time"

	"gorm.io/gorm"
)

// User
type User struct {
	ID        uint           `json:"id,omitempty" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	FirstName string `json:"firstname,omitempty" validate:"required" gorm:"primaryKey"`
	LastName  string `json:"lastname,omitempty" validate:"required"`
	Age       uint8  `json:"age,omitempty" validate:"gte=0,lte=130"`
	Email     string `json:"email,omitempty" validate:"required,email"`
}
