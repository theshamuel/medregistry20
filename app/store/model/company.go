package model

type Company struct {
	ShortName string `json:"shortName"`
	FullName  string `json:"fullName"`
	License   string `json:"license"`
}
