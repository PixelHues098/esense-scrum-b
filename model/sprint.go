package model

import (
	"esense/database"
	"time"

	"gorm.io/gorm"
)

type Sprint struct {
	gorm.Model
	Name        string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	ProjectID   uint
	CreatorID   uint
	IsDone      bool
	Issues      []Issue
}

func (sprint *Sprint) Save() (*Sprint, error) {
	err := database.Database.Create(&sprint).Error
	if err != nil {
		return &Sprint{}, err
	}
	return sprint, nil
}
