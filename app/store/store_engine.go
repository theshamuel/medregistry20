package store

import (
	"github.com/theshamuel/medregistry20/app/store/model"
)

type EngineInterface interface {
	FindVisitByID(id string) (model.Visit, error)
	FindClientByID(id string) (model.Client, error)
	FindDoctorByID(id string) (model.Doctor, error)
	CompanyDetail() (model.Company, error)
	FindVisitsByDoctorSinceTill(doctorID string, startDateEvent, endDateEvent string) ([]model.Visit, error)
	FindVisitsByClientIDSinceTill(clientID string, startDateEventStr, endDateEventStr string) ([]model.Visit, error)
	GetNalogSpravkaSeq() (int, error)
	IncrementNalogSpravkaSeq(idx int) error
	Close() error
}
