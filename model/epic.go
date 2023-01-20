package model

import (
	"esense/database"

	"gorm.io/gorm"
)

type Epic struct {
	gorm.Model
	Name      string
	StartDate string
	EndDate   string
	ProjectID uint
	Issues    []Issue
}

func (epic *Epic) Save() (*Epic, error) {
	err := database.Database.Create(&epic).Error
	if err != nil {
		return &Epic{}, err
	}
	return epic, nil
}
