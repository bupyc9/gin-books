package authors

import (
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
	authors := router.Group("/authors")

	authors.GET("", handler.list)
	authors.POST("", handler.create)
	authors.GET("/:id", handler.find)
	authors.DELETE("/:id", handler.delete)
}

func (handler handler) create(context *gin.Context) {
	var request CreateAuthor
	if err := context.ShouldBind(&request); err != nil {
		context.Error(err)
		return
	}

	author := Author{FirstName: request.FirstName, LastName: request.LastName, SecondName: request.SecondName}
	result := handler.DB.Create(&author)

	if result.Error != nil {
		context.Error(errors.New("failed to create author"))
		return
	}

	context.JSON(http.StatusCreated, author)
}

func (handler handler) find(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))

	var author Author
	result := handler.DB.First(&author, id)
	if result.Error != nil {
		context.Error(result.Error)

		return
	}

	context.JSON(http.StatusOK, author)
}

func (handler handler) delete(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))

	var author Author
	result := handler.DB.First(&author, id)
	if result.Error != nil {
		context.Error(result.Error)

		return
	}

	result = handler.DB.Delete(&author)
	if result.Error != nil {
		context.Error(result.Error)

		return
	}

	context.Status(http.StatusNoContent)
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

	var authors []Author
	result := handler.DB.Limit(limit).Offset(offset).Order("id asc").Find(&authors)
	if result.Error != nil {
		context.Error(result.Error)

		return
	}

	context.JSON(http.StatusOK, authors)
}
