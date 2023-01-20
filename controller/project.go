package controller

import (
	"errors"
	"esense/database"
	"esense/helper"
	"esense/model"
	"esense/model/request"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AddProject(context *gin.Context) {
	var requestBody request.CreateProject
	var usersToAdd []model.User

	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectOwner, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	usersToAdd, err = model.FindUsersByEmail(requestBody.Members)
	if helper.DidContextErr(err, context) {
		return
	}

	var projectToCreate = model.Project{
		Name:        helper.CapitalizeFirstLetter(requestBody.Name),
		Description: requestBody.Description,
		Key:         strings.ToUpper(requestBody.Key),
		Type:        requestBody.Type,
		Members:     append(usersToAdd, projectOwner), /* Add owner to their own project */
		OwnerID:     projectOwner.ID,
	}

	savedProject, err := projectToCreate.Save()
	if helper.DidContextErr(err, context) {
		return
	}

	err = savedProject.CreateBaseSwimlane()
	if helper.DidContextErr(err, context) {
		return
	}

	err = savedProject.CreateBacklog()
	if helper.DidContextErr(err, context) {
		return
	}
}

func GetOwnProjects(context *gin.Context) {
	user, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	err = database.Database.Preload("OwnedProjects").Where("ID=?", user.ID).Find(&user).Error
	if helper.DidContextErr(err, context) {
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": user.OwnedProjects})
}

func GetJoinedProjects(context *gin.Context) {
	user, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	err = database.Database.Preload(clause.Associations).Where("ID=?", user.ID).Find(&user).Error
	if helper.DidContextErr(err, context) {
		return
	}

	context.JSON(http.StatusOK, user.Projects)
}

func AddUserToProject(context *gin.Context) {
	var requestBody request.UserToAddToProject
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userToAdd, err := model.FindUserByEmail(requestBody.Email)
	if helper.DidContextErr(err, context) {
		return
	}

	projectToJoin, err := model.FindProjectById(requestBody.ProjectID)
	if helper.DidContextErr(err, context) {
		return
	}

	database.Database.Model(&userToAdd).Association("Projects").Append([]model.Project{projectToJoin})
	context.JSON(http.StatusOK, gin.H{})
}

func GetProject(context *gin.Context) {
	var project model.Project
	var requestBody request.GetProjectById
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	err = database.Database.Preload(clause.Associations).Preload("Swimlanes", func(db *gorm.DB) *gorm.DB {
		db = db.Order("Position asc")
		return db
	}).Preload("Sprints", func(db *gorm.DB) *gorm.DB {
		db = db.Order("ID asc")
		return db
	}).Where("id=?", requestBody.ID).First(&project).Error
	if helper.DidContextErr(err, context) {
		return
	}

	for _, member := range project.Members {
		if member.ID == user.ID {
			context.JSON(http.StatusOK, project)
			return
		}
	}
	context.JSON(http.StatusBadRequest, gin.H{})
}

func UpdateProjectInfo(context *gin.Context) {
	var requestBody request.UpdateProjectInfo
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prevOwner, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	project, err := model.FindProjectById(requestBody.ProjectID)
	if helper.DidContextErr(err, context) {
		return
	}

	if !project.IsUserMember(prevOwner.ID) {
		return
	}

	if prevOwner.ID != project.OwnerID {
		err = errors.New("user does not own the project!")
		if helper.DidContextErr(err, context) {
			return
		}
	}

	newOwner, err := model.FindUserByEmail(requestBody.NewOwnerEmail)
	err = database.Database.Model(&project).Updates(map[string]interface{}{
		"Name":        requestBody.Name,
		"Description": requestBody.Description,
		"OwnerID":     newOwner.ID,
	}).Error
	if helper.DidContextErr(err, context) {
		return
	}

	context.JSON(http.StatusOK, project)
}
