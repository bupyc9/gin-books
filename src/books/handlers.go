package books

import (
	"books/authors"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func Init(db *gorm.DB, router *gin.RouterGroup) {
	handler := &handler{DB: db}
	handler.Route(router)
}

func (handler handler) Route(router *gin.RouterGroup) {
	authors := router.Group("/books")

	authors.POST("", handler.create)
}

func (handler handler) create(context *gin.Context) {
	var request CreateBookRequest
	if err := context.ShouldBind(&request); err != nil {
		context.Error(err)
		return
	}

	var author authors.Author
	handler.DB.First(&author, request.AuthorID)

	book := Book{
		Author: author,
		Name:   request.Name,
		Pages:  request.Pages,
		Year:   request.Year,
	}

	result := handler.DB.Create(&book)

	if result.Error != nil {
		context.Error(errors.New("failed to create author"))
		return
	}

	context.JSON(http.StatusCreated, book)
}
