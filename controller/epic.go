package controller

import (
	"esense/helper"
	"esense/model"
	"esense/model/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddEpic(context *gin.Context) {
	var requestBody request.CreateEpic
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

	var epicToCreate = model.Epic{
		Name:      requestBody.Name,
		StartDate: requestBody.StartDate,
		EndDate:   requestBody.EndDate,
		ProjectID: requestBody.ProjectID,
	}

	newEpic, err := epicToCreate.Save()
	if helper.DidContextErr(err, context) {
		return
	}

	context.JSON(http.StatusOK, newEpic)
}
