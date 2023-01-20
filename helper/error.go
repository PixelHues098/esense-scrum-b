package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DidContextErr(err error, context *gin.Context) bool {
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return true
	}
	return false
}
