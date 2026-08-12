package books

import (
	"books/authors"
	"time"
)

type Book struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Name      string         `json:"name"`
	Year      uint           `json:"year"`
	Pages     uint           `json:"pages"`
	AuthorID  uint           `json:"-"`
	Author    authors.Author `json:"author" gorm:"constraint:OnDelete:CASCADE;"`
}
