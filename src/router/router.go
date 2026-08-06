package router

import (
	"books/authors"
	"log"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ValidationError struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

type Error struct {
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
}

func errorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		if len(context.Errors) == 0 {
			return
		}
		err := context.Errors.Last().Err

		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			log.Println("validation errors:", validateErrs)
			validationError(validateErrs, context)

			return
		}

		context.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
}

func validationError(validateErrs validator.ValidationErrors, context *gin.Context) {
	response := ValidationError{Message: "Validation Error"}
	response.Errors = make(map[string]string)
	for _, e := range validateErrs {
		response.Errors[e.Field()] = e.Tag()
	}
	context.JSON(http.StatusUnprocessableEntity, response)
}
