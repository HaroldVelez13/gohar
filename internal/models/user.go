package models

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=3,max=150"`
	Age   int8   `json:"age" validate:"required,gt=5,lte=120"`
}