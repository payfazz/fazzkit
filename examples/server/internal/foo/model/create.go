package model

//CreateFoo ...
type CreateFoo struct {
	ID    int    `json:"id" httpurl:"id" validate:"required"`
	Name  string `json:"name" validate:"required,min=3"`
	Price int    `json:"price" validate:"required,min=0"`
}
