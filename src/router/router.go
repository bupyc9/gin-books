package router

import (
	"books/authors"
	"books/books"

	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ValidationResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(errorHandler())
	api(db, router)

	return router
}

func api(db *gorm.DB, router *gin.Engine) {
	api := router.Group("/api")

	authors.Init(db, api)
	books.Init(db, api)
}

func errorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		if len(context.Errors) == 0 {
			return
		}
		err := context.Errors.Last().Err

		var validateErrs validator.ValidationErrors
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {
		case errors.As(err, &validateErrs):
			validationError(validateErrs, context)
		case errors.As(err, &unmarshalTypeError):
			unMarshalError(unmarshalTypeError, context)
		case errors.Is(err, gorm.ErrRecordNotFound):
			recordNotFound(err, context)
		default:
			context.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
		}
	}
}

func validationError(validateErrs validator.ValidationErrors, context *gin.Context) {
	response := ValidationResponse{Message: "Validation Error"}
	response.Errors = make(map[string]string)

	mapValidations := map[string]string{
		"required": "required",
		"author":   "not_found",
	}

	for _, e := range validateErrs {
		value, ok := mapValidations[e.Tag()]
		if !ok {
			value = "unknown"
		}

		response.Errors[e.Field()] = value
	}
	context.JSON(http.StatusUnprocessableEntity, response)
}

func recordNotFound(err error, context *gin.Context) {
	response := MessageResponse{Message: err.Error()}

	context.JSON(http.StatusNotFound, response)
}

func unMarshalError(unmarshalError *json.UnmarshalTypeError, context *gin.Context) {
	response := ValidationResponse{Message: "Validation Error"}
	response.Errors = map[string]string{unmarshalError.Field: fmt.Sprintf("The expected type is \"%s\"", unmarshalError.Type.Name())}
	context.JSON(http.StatusUnprocessableEntity, response)
}
