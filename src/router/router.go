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

		if errors.Is(err, gorm.ErrRecordNotFound) {
			recordNotFound(err, context)

			return
		}

		context.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
	}
}

func validationError(validateErrs validator.ValidationErrors, context *gin.Context) {
	response := ValidationResponse{Message: "Validation Error"}
	response.Errors = make(map[string]string)
	for _, e := range validateErrs {
		response.Errors[e.Field()] = e.Tag()
	}
	context.JSON(http.StatusUnprocessableEntity, response)
}

func recordNotFound(err error, context *gin.Context) {
	response := MessageResponse{Message: err.Error()}

	context.JSON(http.StatusNotFound, response)
}
