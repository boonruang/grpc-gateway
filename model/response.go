package model

import (
	"github.com/gin-gonic/gin"
)

const (
	API_RESPONSE_MESSAGE_400 = "bad request"
	API_RESPONSE_MESSAGE_502 = "bad gateway"
	API_RESPONSE_MESSAGE_504 = "gateway timeout"
)

type APIResponsePayload struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func APIResponse(ctx *gin.Context, statusCode int, message string, data interface{}) APIResponsePayload {
	response := APIResponsePayload{
		Code:    statusCode,
		Message: message,
		Data:    data,
	}

	ctx.AbortWithStatusJSON(statusCode, response)
	return response
}
