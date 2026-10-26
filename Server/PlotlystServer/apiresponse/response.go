package apiresponse

import "github.com/gin-gonic/gin"

type ErrorBody struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func Error(c *gin.Context, status int, code, message string, details ...any) {
	body := ErrorBody{Error: APIError{Code: code, Message: message}}
	if len(details) > 0 {
		body.Error.Details = details[0]
	}
	c.AbortWithStatusJSON(status, body)
}
