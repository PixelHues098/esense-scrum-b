package controller

import (
	"esense/database"
	"esense/helper"
	"esense/model"
	"esense/model/request"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AddSprint(context *gin.Context) {
	var requestBody request.CreateSprint
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sprintCreator, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	project, err := model.FindProjectById(requestBody.ProjectID)
	if helper.DidContextErr(err, context) {
		return
	}

	if !project.IsUserMember(sprintCreator.ID) {
		return
	}

	var sprintToCreate = model.Sprint{
		Name:        helper.CapitalizeFirstLetter(requestBody.Name),
		Description: requestBody.Description,
		IsDone:      false,
		CreatorID:   sprintCreator.ID,
		ProjectID:   requestBody.ProjectID,
	}

	newSprint, err := sprintToCreate.Save()
	if helper.DidContextErr(err, context) {
		return
	}

	context.JSON(http.StatusOK, newSprint)
}

func StartSprint(context *gin.Context) {
	var project model.Project
	var sprint model.Sprint
	var requestBody request.StartSprint
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	projectToCheck, err := model.FindProjectById(requestBody.ProjectID)
	if helper.DidContextErr(err, context) {
		return
	}

	if !projectToCheck.IsUserMember(user.ID) {
		context.JSON(http.StatusBadRequest, gin.H{})
	}

	database.Database.Model(&project).Where("id = ?", requestBody.ProjectID).Update("ActiveSprint", requestBody.SprintID)
	database.Database.Model(&sprint).Where("id = ?", requestBody.SprintID).Update("StartDate", time.Now())
	context.JSON(http.StatusOK, project)
}

func EndSprint(context *gin.Context) {
	var project model.Project
	var sprintToComplete model.Sprint
	var backlog model.Sprint
	var doneSwimlane model.Swimlane
	var requestBody request.EndSprint
	var issuesToUpdate []model.Issue

	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	projectToCheck, err := model.FindProjectById(requestBody.ProjectID)
	if helper.DidContextErr(err, context) {
		return
	}

	if !projectToCheck.IsUserMember(user.ID) {
		context.JSON(http.StatusBadRequest, gin.H{})
	}

	err = database.Database.Where("id = ?", requestBody.ProjectID).First(&project).Update("ActiveSprint", nil).Error
	if helper.DidContextErr(err, context) {
		return
	}

	err = database.Database.Where("id = ?", requestBody.SprintID).First(&sprintToComplete).Updates(model.Sprint{EndDate: time.Now(), IsDone: true}).Error
	if helper.DidContextErr(err, context) {
		return
	}

	err = database.Database.Where("project_id = ? AND name = ?", requestBody.ProjectID, "Done").First(&doneSwimlane).Error
	if helper.DidContextErr(err, context) {
		return
	}

	err = database.Database.Where("project_id = ? AND name = ?", requestBody.ProjectID, "Backlog").First(&backlog).Error
	if helper.DidContextErr(err, context) {
		return
	}

	err = database.Database.Where("sprint_id = ?", requestBody.SprintID).Not("swimlane_id = ?", doneSwimlane.ID).Find(&issuesToUpdate).Update("SprintID", backlog.ID).Error
	if helper.DidContextErr(err, context) {
		return
	}

	context.JSON(http.StatusOK, gin.H{})
}
