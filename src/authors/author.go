package authors

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	FirstName  string         `json:"firstName"`
	LastName   string         `json:"lastName"`
	SecondName string         `json:"secondName"`
}
