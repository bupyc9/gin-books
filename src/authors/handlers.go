package authors

import (
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
	authors := router.Group("/authors")

	authors.POST("", handler.create)
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
