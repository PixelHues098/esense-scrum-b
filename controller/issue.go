package controller

import (
	"esense/database"
	"esense/helper"
	"esense/model"
	"esense/model/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddIssue(context *gin.Context) {
	var requestBody request.CreateIssue
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	issueCreator, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	project, err := model.FindProjectById(requestBody.ProjectID)
	if helper.DidContextErr(err, context) {
		return
	}

	if !project.IsUserMember(issueCreator.ID) {
		return
	}

	todoSwimlane, err := project.GetToDoSwimlane()
	if helper.DidContextErr(err, context) {
		return
	}

	var assigneeId = model.GetValidAssigneeID(requestBody.AssigneeID)
	var epicId = model.GetValidEpicID(requestBody.EpicID)

	var issueToCreate = model.Issue{
		ID:          project.GenerateNewIssueId(),
		Title:       requestBody.Title,
		Description: requestBody.Description,
		Type:        requestBody.Type,
		ProjectID:   project.ID,
		CreatorID:   issueCreator.ID,
		AssigneeID:  assigneeId,
		SprintID:    requestBody.SprintID,
		Points:      requestBody.Points,
		Priority:    requestBody.Priority,
		SwimlaneID:  todoSwimlane.ID,
		EpicID:      epicId,
	}

	newIssue, err := issueToCreate.Save()
	if helper.DidContextErr(err, context) {
		return
	}

	context.JSON(http.StatusOK, newIssue)
}

func UpdateIssue(context *gin.Context) {
	var requestBody request.UpdateIssue
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doer, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	project, err := model.FindProjectById(requestBody.ProjectID)
	if helper.DidContextErr(err, context) {
		return
	}

	if !project.IsUserMember(doer.ID) {
		return
	}

	issueToUpdate, err := model.GetIssueByID(requestBody.ID)
	if helper.DidContextErr(err, context) {
		return
	}

	var assigneeId = model.GetValidAssigneeID(requestBody.AssigneeID)
	var epicId = model.GetValidEpicID(requestBody.EpicID)

	err = database.Database.Model(&issueToUpdate).Updates(map[string]interface{}{
		"Title":       requestBody.Title,
		"Description": requestBody.Description,
		"Type":        requestBody.Type,
		"AssigneeID":  assigneeId,
		"SprintID":    requestBody.SprintID,
		"EpicID":      epicId,
		"Points":      requestBody.Points,
		"Priority":    requestBody.Priority,
	}).Error

	if helper.DidContextErr(err, context) {
		return
	}

	context.JSON(http.StatusOK, issueToUpdate)
}

func DeleteIssue(context *gin.Context) {
	var requestBody request.DeleteIssue
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doer, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	project, err := model.FindProjectById(requestBody.ProjectID)
	if helper.DidContextErr(err, context) {
		return
	}

	if !project.IsUserMember(doer.ID) {
		return
	}

	issueToDelete, err := model.GetIssueByID(requestBody.ID)
	if helper.DidContextErr(err, context) {
		return
	}

	err = database.Database.Delete(&issueToDelete).Error

	if helper.DidContextErr(err, context) {
		return
	}

	context.JSON(http.StatusOK, issueToDelete)
}

func MoveIssueSprint(context *gin.Context) {
	var requestBody request.MoveIssueSprint
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doer, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	issueToMove, err := model.GetIssueByID(requestBody.IssueID)
	if helper.DidContextErr(err, context) {
		return
	}

	project, err := model.FindProjectById(issueToMove.ProjectID)
	if helper.DidContextErr(err, context) {
		return
	}

	if !project.IsUserMember(doer.ID) {
		return
	}

	// TODO: check if sprintId is legit

	database.Database.Model(&issueToMove).Update("SprintID", requestBody.RelocSprintID)
	context.JSON(http.StatusOK, issueToMove)
}

func MoveIssueSwimlane(context *gin.Context) {
	var requestBody request.MoveIssueSwimlane
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doer, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	issueToMove, err := model.GetIssueByID(requestBody.IssueID)
	if helper.DidContextErr(err, context) {
		return
	}

	project, err := model.FindProjectById(issueToMove.ProjectID)
	if helper.DidContextErr(err, context) {
		return
	}

	if !project.IsUserMember(doer.ID) {
		return
	}

	database.Database.Model(&issueToMove).Update("SwimlaneID", requestBody.RelocSwimlaneID)
	context.JSON(http.StatusOK, issueToMove)
}
