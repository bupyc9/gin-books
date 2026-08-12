package books

type CreateBookRequest struct {
	AuthorID uint   `form:"authorId" json:"authorId" binding:"required,author"`
	Name     string `form:"name" json:"name" binding:"required"`
	Pages    uint   `form:"pages" json:"pages" binding:"required"`
	Year     uint   `form:"year" json:"year" binding:"required"`
}
