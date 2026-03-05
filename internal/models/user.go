package models

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=3,max=150"`
	Age   int8   `json:"age" validate:"required,gt=5,lte=120"`
}

// UserPagedResponse estructura la respuesta con metadatos
type UserPagedResponse struct {
	Data     []User `json:"data"`
	Total    int    `json:"total"`
	Page     int    `json:"page"`
	LastPage int    `json:"last_page"`
}
