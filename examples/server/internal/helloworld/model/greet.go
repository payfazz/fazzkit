package model

//Greet ...
type Greet struct {
	Name string `json:"name" httpurl:"name" validate:"required"`
}

//GreetResponse ...
type GreetResponse struct {
	Message string `json:"message"`
}
