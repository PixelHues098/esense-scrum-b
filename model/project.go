package model

import (
	"esense/database"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Project struct {
	gorm.Model
	Name         string `gorm:"type:text" json:"name"`
	Description  string `gorm:"type:text" json:"desc"`
	Key          string `gorm:"type:text; unique" json:"key"`
	Type         string `gorm:"type:text" json:"type"`
	OwnerID      uint   `gorm:"type:text" json:"ownerId"`
	Members      []User `gorm:"many2many:user_projects;"`
	Issues       []Issue
	Sprints      []Sprint
	Swimlanes    []Swimlane
	ActiveSprint uint
}

func (project *Project) Save() (*Project, error) {
	err := database.Database.Create(&project).Error
	if err != nil {
		return &Project{}, err
	}
	return project, nil
}

func FindProjectById(id uint) (Project, error) {
	var project Project
	err := database.Database.Preload(clause.Associations).Where("ID=?", id).Find(&project).Error
	if err != nil {
		return Project{}, err
	}

	return project, nil
}

func (project *Project) IsUserMember(idToCheck uint) bool {
	for _, member := range project.Members {
		if member.ID == idToCheck {
			return true
		}
	}
	return false
}

func (project *Project) GenerateNewIssueId() string {
	projectKey := project.Key
	idOffset := 1

	projectIssueCnt := len(project.Issues) + idOffset
	projectIssueCntStr := strconv.Itoa(projectIssueCnt)

	newIssueId := projectKey + "-" + projectIssueCntStr

	for project.CheckIfIssueIdDuplicate(newIssueId) {
		idOffset++

		projectIssueCnt = len(project.Issues) + idOffset
		projectIssueCntStr = strconv.Itoa(projectIssueCnt)

		newIssueId = projectKey + "-" + projectIssueCntStr
	}

	return newIssueId
}

func (project *Project) CheckIfIssueIdDuplicate(issueId string) bool {
	for _, issue := range project.Issues {
		if issue.ID == issueId {
			return true
		}
	}

	return false
}

func (project *Project) CreateBaseSwimlane() error {

	var baseSwimlane = []string{"To Do", "In Progress", "Done"}

	for index, swimlaneName := range baseSwimlane {
		var position = uint(index)
		if swimlaneName == "Done" {
			position = 999
		}
		var swimlane = Swimlane{
			Name:      swimlaneName,
			Position:  position,
			ProjectID: project.ID,
		}

		_, err := swimlane.Save()
		if err != nil {
			return err
		}
	}
	return nil
}

func (project *Project) CreateBacklog() error {

	var backlog = Sprint{
		Name:        "Backlog",
		Description: "",
		ProjectID:   project.ID,
		CreatorID:   project.OwnerID,
	}

	_, err := backlog.Save()
	if err != nil {
		return err
	}

	return nil
}
