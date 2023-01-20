package controller

import (
	"esense/helper"
	"esense/model"
	"esense/model/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddSwimlane(context *gin.Context) {
	var requestBody request.CreateSwimlane
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	swimlaneCreator, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	project, err := model.FindProjectById(requestBody.ProjectID)
	if helper.DidContextErr(err, context) {
		return
	}

	if !project.IsUserMember(swimlaneCreator.ID) {
		return
	}

	var swimlaneToCreate = model.Swimlane{
		Name:      helper.CapitalizeFirstLetter(requestBody.Name),
		Position:  uint(len(project.Swimlanes) - 1),
		ProjectID: requestBody.ProjectID,
	}

	newSwimlane, err := swimlaneToCreate.Save()
	if helper.DidContextErr(err, context) {
		return
	}

	context.JSON(http.StatusOK, newSwimlane)
}
