package store

import (
	"github.com/theshamuel/medregistry20/app/store/model"
)

type EngineInterface interface {
	FindVisitByID(id string) (model.Visit, error)
	FindClientByID(id string) (model.Client, error)
	FindDoctorByID(id string) (model.Doctor, error)
	FindDoctors() (doctors []model.Doctor, err error)
	CompanyDetail() (model.Company, error)
	FindVisitsByDoctorSinceTill(doctorID string, startDateEvent, endDateEvent string) ([]model.Visit, error)
	FindVisitsByClientIDSinceTill(clientID string, startDateEventStr, endDateEventStr string) ([]model.Visit, error)
	GetProfitByDoctorSinceTill(startDateEventStr, endDateEventStr string) ([]model.ProfitByDoctorSinceTillRecord, error)
	IncrementSeq(idx int, code string) error
	GetSeq(code string) (int, error)
	Close() error
}
