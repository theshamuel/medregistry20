package model

type Doctor struct {
	ID               string `json:"id"`
	FirstName        string `json:"name"`
	Surname          string `json:"surname"`
	Middlename       string `json:"middlename"`
	PositionGenitive string `json:"positionLabelGenitive"` //Genitive form for report field
}
