package books

import (
	"books/authors"
	"errors"
	"net/http"
	"strconv"

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
	authors.GET("/:id", handler.find)
	authors.DELETE("/:id", handler.delete)
	authors.GET("", handler.list)
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

func (handler handler) find(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))

	var book Book
	result := handler.DB.Preload("Author").First(&book, id)
	if result.Error != nil {
		context.Error(result.Error)

		return
	}

	context.JSON(http.StatusOK, book)
}

func (handler handler) delete(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))

	var book Book
	result := handler.DB.First(&book, id)
	if result.Error != nil {
		context.Error(result.Error)

		return
	}

	result = handler.DB.Delete(&book)
	if result.Error != nil {
		context.Error(result.Error)

		return
	}

	context.JSON(http.StatusOK, book)
}

func (handler handler) list(context *gin.Context) {
	limit, _ := strconv.Atoi(context.DefaultQuery("count", "40"))
	if limit <= 0 {
		limit = 40
	}
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	var books []Book
	result := handler.DB.Limit(limit).Offset(offset).Preload("Author").Find(&books)
	if result.Error != nil {
		context.Error(result.Error)

		return
	}

	context.JSON(http.StatusOK, books)
}
