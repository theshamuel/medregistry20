package model

import "time"

type Visit struct {
	ID                        string    `json:"id"`
	DoctorID                  string    `json:"doctorId"`
	DoctorName                string    `json:"doctorLabel"`
	DoctorExcludedFromReports bool      `json:"doctorExcludedFromReports"`
	ClientID                  string    `json:"clientId"`
	TotalSum                  int       `json:"totalSum"`
	Services                  []Service `json:"services"`
	DateEventStr              string    `json:"dateEvent"`
	DateEvent                 time.Time
	Diagnosis                 string `json:"diagnosis"`
	AdditionalExamination     string `json:"additionalExamination"`
	Therapy                   string `json:"therapy"`
}
