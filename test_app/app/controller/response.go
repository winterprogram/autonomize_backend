package controller

import (
	responsedto "test/test_app/app/model/response"

	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, responsedto.ErrorResponse{Success: false, Error: responsedto.ErrorResponseData{Code: code, Message: message}})
}

func ResponseMethod(c *gin.Context, code int, body interface{}) {
	c.JSON(code, body)
}
