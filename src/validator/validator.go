package validator

import (
	"books/authors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Validator struct {
	DB *gorm.DB
}

func Init(db *gorm.DB) {
	v := Validator{DB: db}

	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		val.RegisterValidation("author", v.author)
	}
}

func (v *Validator) author(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(uint)
	if !ok {
		return false
	}

	var author authors.Author
	result := v.DB.First(&author, value)
	if result.Error != nil {
		return false
	}

	return true
}
