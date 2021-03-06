package store

import (
	"github.com/theshamuel/medregistry20/app/store/model"
)

type EngineInterface interface {
	FindVisitByID(id string)(model.Visit, error)
	FindClientByID(id string)(model.Client, error)
	FindDoctorByID(id string)(model.Doctor, error)
	CompanyDetail()(model.Company, error)
	FindVisitByDoctorSinceTill(doctorID string, startDateEvent, endDateEvent string)([]model.Visit, error)
	Close() error
}
