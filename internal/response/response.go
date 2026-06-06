package response

import (
	"errors"
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
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Details    any
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(statusCode int, code string, message string, err error, details any) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    details,
		Err:        err,
	}
}

func NewBadRequest(message string, err error, details any) *AppError {
	return NewAppError(http.StatusBadRequest, CodeBadRequest, message, err, details)
}

func NewUnauthorized(message string, err error, details any) *AppError {
	return NewAppError(http.StatusUnauthorized, CodeUnauthorized, message, err, details)
}

func NewForbidden(message string, err error, details any) *AppError {
	return NewAppError(http.StatusForbidden, CodeForbidden, message, err, details)
}

func NewNotFound(message string, err error, details any) *AppError {
	return NewAppError(http.StatusNotFound, CodeNotFound, message, err, details)
}

func NewConflict(message string, err error, details any) *AppError {
	return NewAppError(http.StatusConflict, CodeConflict, message, err, details)
}

func NewInternalServerError(message string, err error, details any) *AppError {
	return NewAppError(http.StatusInternalServerError, CodeInternalServerError, message, err, details)
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
			Message: message,
			Details: details,
		},
	})
}

func AppErrorResponse(c *gin.Context, appErr *AppError) {
	if appErr == nil {
		InternalServerError(c, "internal server error", nil)
		return
	}

	Error(c, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
}

func HandleError(c *gin.Context, err error) {
	if err == nil {
		AppErrorResponse(c, NewInternalServerError("internal server error", nil, nil))
		return
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		AppErrorResponse(c, appErr)
		return
	}

	AppErrorResponse(c, NewInternalServerError("internal server error", err, nil))
}

func BadRequest(c *gin.Context, message string, details any) {
	AppErrorResponse(c, NewBadRequest(message, nil, details))
}

func Unauthorized(c *gin.Context, message string, details any) {
	AppErrorResponse(c, NewUnauthorized(message, nil, details))
}

func Forbidden(c *gin.Context, message string, details any) {
	AppErrorResponse(c, NewForbidden(message, nil, details))
}

func NotFound(c *gin.Context, message string, details any) {
	AppErrorResponse(c, NewNotFound(message, nil, details))
}

func Conflict(c *gin.Context, message string, details any) {
	AppErrorResponse(c, NewConflict(message, nil, details))
}

func InternalServerError(c *gin.Context, message string, details any) {
	AppErrorResponse(c, NewInternalServerError(message, nil, details))
}
