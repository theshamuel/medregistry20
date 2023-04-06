package store

import (
	"encoding/json"
	"github.com/theshamuel/medregistry20/app/store/model"
	"github.com/theshamuel/medregistry20/app/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Mix struct {
	URI         string
	HttpClient  *utils.Repeater
	MongoClient *mongo.Client
}

func (s *Mix) FindVisitByDoctorSinceTill(doctorID string, startDateEvent, endDateEvent string) ([]model.Visit, error) {
	log.Printf("[INFO] FindVisitByDoctorSinceTill param doctorID=%s;startDateEvent=%s;endDateEvent=%s;",
		doctorID, startDateEvent, endDateEvent)
	s.HttpClient = &utils.Repeater{
		ClientTimeout: 10,
		Attempts:      10,
		URI:           s.URI + "/visits/" + doctorID + "/" + startDateEvent + "/" + endDateEvent + "/",
		Count:         3,
	}
	data, err := s.HttpClient.Get()
	if err != nil {
		log.Printf("[ERROR] cannot receive data from MedRegistry API v1")
	}
	var visits []model.Visit
	err = json.Unmarshal(data, &visits)
	if err != nil {
		log.Printf("[ERROR] cannot unmarshar response %#v", err)
		return nil, err
	}
	return visits, nil
}

func (s *Mix) FindVisitByID(visitID string) (visit model.Visit, err error) {
	log.Printf("[INFO] FindVisitByID params visitID=%s;", visitID)
	s.HttpClient = &utils.Repeater{
		ClientTimeout: 10,
		Attempts:      10,
		URI:           s.URI + "/visits/" + visitID,
		Count:         3,
	}
	data, err := s.HttpClient.Get()
	if err != nil {
		log.Printf("[ERROR] cannot receive data from MedRegistry API v1")
	}

	err = json.Unmarshal(data, &visit)

	if visit.DateEventStr != "" {
		visit.DateEvent, _ = time.Parse("2006-01-02", visit.DateEventStr)
	}

	if err != nil {
		log.Printf("[ERROR] cannot unmarshar response %#v", err)
		return visit, err
	}
	return visit, nil
}
func (s *Mix) FindClientByID(clientID string) (client model.Client, err error) {
	log.Printf("[INFO] FindClientByID params clientID=%s;", clientID)
	s.HttpClient = &utils.Repeater{
		ClientTimeout: 10,
		Attempts:      10,
		URI:           s.URI + "/clients/" + clientID,
		Count:         3,
	}
	data, err := s.HttpClient.Get()
	if err != nil {
		log.Printf("[ERROR] cannot receive data from MedRegistry API v1")
	}

	err = json.Unmarshal(data, &client)

	if client.BirthdayStr != "" {
		client.Birthday, _ = time.Parse("2006-01-02", client.BirthdayStr)
		client.Age = int((time.Since(client.Birthday).Hours() / 24) / 365)
		client.SetAgePostfix()
	}

	if err != nil {
		log.Printf("[ERROR] cannot unmarshar response %#v", err)
		return client, err
	}
	return client, nil
}
func (s *Mix) FindDoctorByID(doctorID string) (doctor model.Doctor, err error) {
	log.Printf("[INFO] FindDoctorByID params doctorID=%s;", doctorID)
	s.HttpClient = &utils.Repeater{
		ClientTimeout: 10,
		Attempts:      10,
		URI:           s.URI + "/doctors/" + doctorID,
		Count:         3,
	}
	data, err := s.HttpClient.Get()
	if err != nil {
		log.Printf("[ERROR] cannot receive data from MedRegistry API v1")
	}

	err = json.Unmarshal(data, &doctor)
	if err != nil {
		log.Printf("[ERROR] cannot unmarshar response %#v", err)
		return doctor, err
	}
	return doctor, nil
}
func (s *Mix) CompanyDetail() (company model.Company, err error) {
	log.Printf("[INFO] get company detail")
	s.HttpClient = &utils.Repeater{
		ClientTimeout: 10,
		Attempts:      10,
		URI:           s.URI + "/company",
		Count:         3,
	}
	data, err := s.HttpClient.Get()
	if err != nil {
		log.Printf("[ERROR] cannot receive data from MedRegistry API v1")
	}

	err = json.Unmarshal(data, &company)
	if err != nil {
		log.Printf("[ERROR] cannot unmarshar response %#v", err)
		return company, err
	}
	return company, nil
}

func (s *Mix) Close() error {
	log.Printf("[INFO] Close Mix")
	return nil
}
