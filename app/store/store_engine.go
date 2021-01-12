package store

import (
	"github.com/theshamuel/medregistry20/app/store/model"
)

type EngineInterface interface {
	FindVisitById(id string)(model.Visit, error)
	FindClientById(id string)(model.Client, error)
	FindDoctorById(id string)(model.Doctor, error)
	CompanyDetail()(model.Company, error)
	FindVisitByDoctorSinceTill(doctorID string, startDateEvent, endDateEvent string)([]model.Visit, error)
	Close() error
}
