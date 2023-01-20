package model

import (
	"errors"
	"esense/database"

	"gorm.io/gorm"
)

type Issue struct {
	gorm.Model
	ID          string `gorm:"primaryKey;autoIncrement:false"`
	Title       string
	Description string
	Type        string
	Priority    string
	Points      uint
	AssigneeID  *uint
	CreatorID   uint
	ProjectID   uint
	SwimlaneID  uint
	SprintID    uint
	EpicID      *uint
}

func (issue *Issue) Save() (*Issue, error) {
	err := database.Database.Create(&issue).Error
	if err != nil {
		return &Issue{}, err
	}
	return issue, nil
}

func GetValidAssigneeID(assigneeId uint) *uint {
	if assigneeId == 0 {
		return nil
	}
	return &assigneeId
}

func (project *Project) GetToDoSwimlane() (Swimlane, error) {
	for _, swimlane := range project.Swimlanes {
		if swimlane.Name == "To Do" {
			return swimlane, nil
		}
	}
	return Swimlane{}, errors.New("todo swimlane was not found in the project")
}

func GetIssueByID(issueId string) (Issue, error) {
	var issue Issue
	err := database.Database.Where("ID=?", issueId).First(&issue).Error

	if err != nil {
		return Issue{}, err
	}

	return issue, nil
}
