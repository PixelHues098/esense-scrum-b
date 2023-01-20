package model

import (
	"esense/database"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName     string    `gorm:"size:255;not null" json:"firstName"`
	LastName      string    `gorm:"size:255;not null" json:"lastName"`
	Username      string    `gorm:"size:255;not null" json:"username"`
	Email         string    `gorm:"size:255;not null;unique" json:"email"`
	Password      string    `gorm:"size:255;not null;" json:"-"`
	Projects      []Project `gorm:"many2many:user_projects;"`
	OwnedProjects []Project `gorm:"foreignKey:OwnerID"`
	AssignedIssue []Issue   `gorm:"foreignKey:AssigneeID"`
	CreatedIssue  []Issue   `gorm:"foreignKey:CreatorID"`
	Sprints       []Sprint  `gorm:"foreignKey:CreatorID"`
}

// BeforeCreate function is a GORM hook: https://gorm.io/docs/hooks.html
func (user *User) BeforeCreate(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(passwordHash)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	return nil
}

func (user *User) Create() (*User, error) {
	err := database.Database.Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func FindUserByEmail(email string) (User, error) {
	var user User
	err := database.Database.Where("email=?", email).First(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func FindUserById(id uint) (User, error) {
	var user User
	err := database.Database.Where("ID=?", id).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func FindUsersByEmail(emails []string) ([]User, error) {
	var users []User
	err := database.Database.Find(&users, "email in ?", emails).Error
	if err != nil {
		return []User{}, err
	}

	return users, nil
}
