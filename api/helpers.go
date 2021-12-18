package api

import "github.com/gin-gonic/gin"

func ErrorResponse(err error) map[string]interface{} {
	return gin.H{
		"error": err.Error(),
	}
}

func SuccessResponse() map[string]interface{} {
	return gin.H{
		"message": "ok",
	}
}
