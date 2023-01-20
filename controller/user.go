package controller

import (
	"esense/database"
	"esense/helper"
	"esense/model"
	"esense/model/request"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

func GetUser(context *gin.Context) {
	var loggedUser model.User
	user, err := helper.GetAuthorizedUserByJWT(context)

	err = database.Database.Preload(clause.Associations).Where("ID=?", user.ID).Find(&loggedUser).Error

	if helper.DidContextErr(err, context) {
		return
	}

	context.JSON(http.StatusOK, loggedUser)
}

func UpdateUserInfo(context *gin.Context) {
	var requestBody request.UpdateUserInfo
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userToUpdate, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	err = database.Database.Model(&userToUpdate).Updates(map[string]interface{}{
		"FirstName": requestBody.FirstName,
		"LastName":  requestBody.LastName,
		"Username":  requestBody.Username,
		"Email":     requestBody.Email,
	}).Error
	if helper.DidContextErr(err, context) {
		return
	}

	context.JSON(http.StatusOK, userToUpdate)
}

func ChangePassword(context *gin.Context) {
	var requestBody request.ChangePassword
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := helper.GetAuthorizedUserByJWT(context)
	if helper.DidContextErr(err, context) {
		return
	}

	err = user.ValidatePassword(requestBody.CurrentPassword)
	if helper.DidContextErr(err, context) {
		return
	}

	err = requestBody.ConfirmPassword()
	if helper.DidContextErr(err, context) {
		return
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(requestBody.NewPassword), bcrypt.DefaultCost)
	passHashStr := string(passHash)
	err = database.Database.Model(&user).Update("Password", passHashStr).Error
	if helper.DidContextErr(err, context) {
		return
	}

	context.JSON(http.StatusOK, user)
}
