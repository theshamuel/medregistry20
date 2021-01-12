package model

type Service struct {
	ID            string         `json:"id"`
	Name          string         `json:"label"`
	Category      string         `json:"category"`
	Price         int            `json:"price"`
	Discount      int            `json:"discount"`
	DoctorPay     int            `json:"doctorPay"`
	DoctorPayType string         `json:"doctorPayType"`
	DoctorID      int            `json:"doctorId"`
	PersonalRates []PersonalRate `json:"personalRates"`
}

type PersonalRate struct {
	ID            string `json:"id"`
	Price         int    `json:"price"`
	DoctorPay     int    `json:"doctorPay"`
	DoctorPayType string `json:"doctorPayType"`
	DoctorID      string `json:"doctorId"`
	DoctorSalary  int    // calculated field = doctorPay * count of service
}
