package exceptions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(message string) interface{} {
	return gin.H{"statusCode": http.StatusBadRequest, "message": message}
}

func BadValidation(errors []ApiError) interface{} {
	return gin.H{"statusCode": http.StatusBadRequest, "errors": errors}
}

func InternalServerError(message string) interface{} {
	return gin.H{"statusCode": http.StatusInternalServerError, "message": message}
}

func NotFound(message string) interface{} {
	return gin.H{"statusCode": http.StatusNotFound, "message": message}
}

func Unauthorized(message string) interface{} {
	return gin.H{"statusCode": http.StatusUnauthorized, "message": message}
}
