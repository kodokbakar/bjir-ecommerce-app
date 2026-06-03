package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeBadRequest          = "bad_request"
	CodeUnauthorized        = "unauthorized"
	CodeForbidden           = "forbidden"
	CodeNotFound            = "not_found"
	CodeConflict            = "conflict"
	CodeInternalServerError = "internal_server_error"
)

type Body struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
	Meta    any        `json:"meta,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Details any    `json:"details,omitempty"`
}

func Success(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, Body{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, statusCode int, message string, data any, meta any) {
	c.JSON(statusCode, Body{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c *gin.Context, statusCode int, code string, message string, details any) {
	c.AbortWithStatusJSON(statusCode, Body{
		Success: false,
		Message: message,
		Error: &ErrorBody{
			Code:    code,
			Details: details,
		},
	})
}

func BadRequest(c *gin.Context, message string, details any) {
	Error(c, http.StatusBadRequest, CodeBadRequest, message, details)
}

func Unauthorized(c *gin.Context, message string, details any) {
	Error(c, http.StatusUnauthorized, CodeUnauthorized, message, details)
}

func Forbidden(c *gin.Context, message string, details any) {
	Error(c, http.StatusForbidden, CodeForbidden, message, details)
}

func NotFound(c *gin.Context, message string, details any) {
	Error(c, http.StatusNotFound, CodeNotFound, message, details)
}

func Conflict(c *gin.Context, message string, details any) {
	Error(c, http.StatusConflict, CodeConflict, message, details)
}

func InternalServerError(c *gin.Context, message string, details any) {
	Error(c, http.StatusInternalServerError, CodeInternalServerError, message, details)
}
