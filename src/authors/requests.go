package authors

type CreateAuthor struct {
	FirstName  string `form:"firstName" json:"firstName" binding:"required"`
	LastName   string `form:"lastName" json:"lastName" binding:"required"`
	SecondName string `form:"secondName" json:"secondName"`
}
