package controller

import (
	"esense/helper"
	"esense/model"
	"esense/model/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(context *gin.Context) {
	var requestBody model.RegisterInput

	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := requestBody.ConfirmPassword(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		FirstName: helper.CapitalizeFirstLetter(requestBody.FirstName),
		LastName:  helper.CapitalizeFirstLetter(requestBody.LastName),
		Username:  requestBody.Username,
		Email:     requestBody.Email,
		Password:  requestBody.Password,
	}

	_, err := user.Create()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{})
}

func Login(context *gin.Context) {
	var requestBody request.Login

	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userToLogin, err := model.FindUserByEmail(requestBody.Email)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = userToLogin.ValidatePassword(requestBody.Password)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := helper.GenerateJWT(userToLogin)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"jwt": jwt})
}
