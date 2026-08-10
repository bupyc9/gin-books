package authors

import (
	"time"
)

type Author struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	SecondName string    `json:"secondName"`
}
