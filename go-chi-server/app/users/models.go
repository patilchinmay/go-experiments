package user

import "gorm.io/gorm"

// User
type User struct {
	gorm.Model
	FirstName string `json:"firstname" validate:"required" gorm:"primaryKey"`
	LastName  string `json:"lastname" validate:"required"`
	Age       uint8  `json:"age" validate:"gte=0,lte=130"`
	Email     string `json:"email" validate:"required,email"`
}
