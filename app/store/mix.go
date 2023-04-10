package store

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/theshamuel/medregistry20/app/store/model"
	"github.com/theshamuel/medregistry20/app/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strings"
	"time"
)

type Mix struct {
	URI         string
	HttpClient  *utils.Repeater
	MongoClient *mongo.Client
}

func (s *Mix) FindVisitsByDoctorSinceTill(doctorID string, startDateEvent, endDateEvent string) ([]model.Visit, error) {
	log.Printf("[INFO] FindVisitsByDoctorSinceTill param doctorID=%s;startDateEvent=%s;endDateEvent=%s;",
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

func (s *Mix) FindVisitsByClientIDSinceTill(clientID string, startDateEventStr, endDateEventStr string) ([]model.Visit, error) {

	res := make([]model.Visit, 0)

	clientIDObj, _ := primitive.ObjectIDFromHex(clientID)
	startDateEvent, err := time.Parse("2006-01-02", startDateEventStr)
	if err != nil {
		return nil, err
	}

	endDateEvent, err := time.Parse("2006-01-02", endDateEventStr)
	if err != nil {
		return nil, err
	}

	startDateEventB := primitive.NewDateTimeFromTime(time.Date(startDateEvent.Year(), startDateEvent.Month(),
		startDateEvent.Day(), 0, 0, 0, 0, time.UTC))
	endDateEventB := primitive.NewDateTimeFromTime(time.Date(endDateEvent.Year(), endDateEvent.Month(),
		endDateEvent.Day(), 23, 59, 59, 0, time.UTC))

	var filter = bson.D{
		{"$and", []bson.D{
			{{"client._id", clientIDObj}},
			{{"dateEvent", bson.D{{"$gte", startDateEventB}}}},
			{{"dateEvent", bson.D{{"$lte", endDateEventB}}}},
		}}}

	visitsCollection := s.MongoClient.Database("medregDB").Collection("visits")
	cursor, err := visitsCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var visitsBson []VisitModel
	if err = cursor.All(context.Background(), &visitsBson); err != nil {
		return nil, err
	}

	for _, vb := range visitsBson {
		res = append(res, model.Visit{
			ID:               vb.ID.Hex(),
			ClientID:         vb.Client.ID.Hex(),
			TotalSum:         vb.CalculateTotalSum(),
			DateEvent:        vb.DateEvent.Time(),
			ClientName:       strings.Title(strings.ToLower(vb.Client.Name)),
			ClientSurname:    strings.Title(strings.ToLower(vb.Client.Surname)),
			ClientMiddlename: strings.Title(strings.ToLower(vb.Client.Middlename)),
		})

	}

	return res, nil
}

func (s *Mix) GetNalogSpravkaSeq() (int, error) {
	filter := bson.M{"code": "nalogSpravkaNum"}
	sequenceCollection := s.MongoClient.Database("medregDB").Collection("sequence")

	var seq SequenceModel
	if err := sequenceCollection.FindOne(context.Background(), filter).Decode(&seq); err != nil {
		return -1, err
	}

	return seq.Seq, nil
}

func (s *Mix) IncrementNalogSpravkaSeq(idx int) error {
	return errors.New("not implemented yet")
}

func (s *Mix) Close() error {
	s.MongoClient.Disconnect(context.Background())
	return nil
}
