package model

import (
	"esense/database"

	"gorm.io/gorm"
)

type Swimlane struct {
	gorm.Model
	Name      string
	Position  uint
	ProjectID uint
	Issues    []Issue
}

func (swimlane *Swimlane) Save() (*Swimlane, error) {
	err := database.Database.Create(&swimlane).Error
	if err != nil {
		return &Swimlane{}, err
	}
	return swimlane, nil
}
