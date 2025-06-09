package dto

type ProductURI struct {
	ID uint `uri:"id" binding:"required,gt=0"`
}
