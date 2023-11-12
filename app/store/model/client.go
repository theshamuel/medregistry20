package model

import "time"

type Client struct {
	ID                string `json:"id"`
	Firstname         string `json:"name"`
	Surname           string `json:"surname"`
	Middlename        string `json:"middlename"`
	BirthdayStr       string `json:"birthday"`
	Sex               string `json:"gender"`
	Email             string `json:"email,omitempty"`
	Phone             string `json:"phone,omitempty"`
	WorkPlace         string `json:"workPlace,omitempty"`
	Occupation        string `json:"workPosition,omitempty"`
	PassportSerial    string `json:"passportSerial"`
	PassportNumber    string `json:"passportNumber"`
	PassportPlace     string `json:"passportPlace,omitempty"`
	PassportCodePlace string `json:"passportCodePlace,omitempty"`
	Address           string `json:"address,omitempty"`
	Birthday          time.Time
	Age               int
	AgePostfix        string
}

func (c *Client) SetAgePostfix() {
	if c.Age > 0 {
		if c.Age > 4 && c.Age < 21 || c.Age >= 10 && c.Age%(c.Age/10*10) == 0 || c.Age >= 10 && c.Age%(c.Age/10*10) > 4 {
			c.AgePostfix = "лет"
		} else if c.Age == 1 || c.Age >= 10 && c.Age%(c.Age/10*10) == 1 {
			c.AgePostfix = "год"
		} else {
			c.AgePostfix = "года"
		}
	}
}
