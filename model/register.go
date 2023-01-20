package model

import (
	"errors"

	"gorm.io/gorm"
)

type RegisterInput struct {
	gorm.Model
	FirstName    string `gorm:"size:255;not null" json:"firstName" binding:"required"`
	LastName     string `gorm:"size:255;not null" json:"lastName" binding:"required"`
	Username     string `gorm:"size:255;not null" json:"username" binding:"required"`
	Email        string `gorm:"size:255;not null" json:"email" binding:"required"`
	Password     string `gorm:"size:255;not null;" json:"password" binding:"required"`
	ConfPassword string `gorm:"size:255;not null;" json:"confPassword" binding:"required"`
}

func (registerInput *RegisterInput) ConfirmPassword() error {
	if registerInput.Password != registerInput.ConfPassword {
		return errors.New("initial password did not match with confirmation password")
	}
	return nil
}
