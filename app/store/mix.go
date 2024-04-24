package store

import (
	"context"
	"encoding/json"
	"github.com/theshamuel/medregistry20/app/store/model"
	"github.com/theshamuel/medregistry20/app/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"strings"
	"time"
)

type Mix struct {
	URI         string
	HTTPClient  *utils.Repeater
	MongoClient *mongo.Client
}

func (s *Mix) FindVisitsByDoctorSinceTill(doctorID string, startDateEvent, endDateEvent string) ([]model.Visit, error) {
	log.Printf("[INFO] FindVisitsByDoctorSinceTill param doctorID=%s;startDateEvent=%s;endDateEvent=%s;",
		doctorID, startDateEvent, endDateEvent)
	s.HTTPClient = &utils.Repeater{
		ClientTimeout: 10,
		Attempts:      10,
		URI:           s.URI + "/visits/" + doctorID + "/" + startDateEvent + "/" + endDateEvent + "/",
		Count:         3,
	}
	data, err := s.HTTPClient.Get()
	if err != nil {
		log.Printf("[ERROR] cannot receive data from MedRegistry API v1")
	}
	var visits []model.Visit
	err = json.Unmarshal(data, &visits)
	if err != nil {
		log.Printf("[ERROR] cannot unmarshal response %#v", err)
		return nil, err
	}
	return visits, nil
}

func (s *Mix) FindVisitByID(visitID string) (visit model.Visit, err error) {
	log.Printf("[INFO] FindVisitByID params visitID=%s;", visitID)
	s.HTTPClient = &utils.Repeater{
		ClientTimeout: 10,
		Attempts:      10,
		URI:           s.URI + "/visits/" + visitID,
		Count:         3,
	}
	data, err := s.HTTPClient.Get()
	if err != nil {
		log.Printf("[ERROR] cannot receive data from MedRegistry API v1")
	}

	err = json.Unmarshal(data, &visit)

	if visit.DateEventStr != "" {
		visit.DateEvent, _ = time.Parse("2006-01-02", visit.DateEventStr)
	}

	if err != nil {
		log.Printf("[ERROR] cannot unmarshal response %#v", err)
		return visit, err
	}
	return visit, nil
}

func (s *Mix) FindClientByID(clientID string) (client model.Client, err error) {
	log.Printf("[INFO] FindClientByID params clientID=%s;", clientID)
	s.HTTPClient = &utils.Repeater{
		ClientTimeout: 10,
		Attempts:      10,
		URI:           s.URI + "/clients/" + clientID,
		Count:         3,
	}
	data, err := s.HTTPClient.Get()
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
		log.Printf("[ERROR] cannot unmarshal response %#v", err)
		return client, err
	}
	return client, nil
}

func (s *Mix) FindDoctorByID(doctorID string) (doctor model.Doctor, err error) {
	log.Printf("[INFO] FindDoctorByID params doctorID=%s;", doctorID)
	s.HTTPClient = &utils.Repeater{
		ClientTimeout: 10,
		Attempts:      10,
		URI:           s.URI + "/doctors/" + doctorID,
		Count:         3,
	}
	data, err := s.HTTPClient.Get()
	if err != nil {
		log.Printf("[ERROR] cannot receive data from MedRegistry API v1")
	}

	err = json.Unmarshal(data, &doctor)
	if err != nil {
		log.Printf("[ERROR] cannot unmarshal response %#v", err)
		return doctor, err
	}
	return doctor, nil
}

func (s *Mix) FindDoctors() (doctors []model.Doctor, err error) {
	log.Println("[INFO] FindDoctors")
	s.HTTPClient = &utils.Repeater{
		ClientTimeout: 10,
		Attempts:      10,
		URI:           s.URI + "/doctors/",
		Count:         3,
	}
	data, err := s.HTTPClient.Get()
	if err != nil {
		log.Printf("[ERROR] cannot receive data from MedRegistry API v1")
	}

	err = json.Unmarshal(data, &doctors)
	if err != nil {
		log.Printf("[ERROR] cannot unmarshal response %#v", err)
		return doctors, err
	}
	return doctors, nil
}

func (s *Mix) CompanyDetail() (company model.Company, err error) {
	log.Printf("[INFO] get company detail")
	s.HTTPClient = &utils.Repeater{
		ClientTimeout: 10,
		Attempts:      10,
		URI:           s.URI + "/company",
		Count:         3,
	}
	data, err := s.HTTPClient.Get()
	if err != nil {
		log.Printf("[ERROR] cannot receive data from MedRegistry API v1")
	}

	err = json.Unmarshal(data, &company)
	if err != nil {
		log.Printf("[ERROR] cannot unmarshal response %#v", err)
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

	caser := cases.Title(language.Russian)

	for _, vb := range visitsBson {
		res = append(res, model.Visit{
			ID:                vb.ID.Hex(),
			ClientID:          vb.Client.ID.Hex(),
			DateEvent:         vb.DateEvent.Time(),
			ClientName:        caser.String(strings.ToLower(vb.Client.Name)),
			ClientSurname:     caser.String(strings.ToLower(vb.Client.Surname)),
			ClientMiddlename:  caser.String(strings.ToLower(vb.Client.Middlename)),
			ClientGender:      vb.Client.Gender,
			TotalSumWithPenny: vb.CalculateTotalSum(),
		})

	}

	return res, nil
}

func (s *Mix) GetProfitByDoctorSinceTill(startDateEventStr, endDateEventStr string) ([]model.ProfitByDoctorSinceTillRecord, error) {

	res := make([]model.ProfitByDoctorSinceTillRecord, 0)

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
	//[
	//    {
	//        $match:
	//            {
	//                'dateEvent':
	//                    {
	//                        $gte: new Date("2023-01-01"),
	//                        $lte: new Date("2023-12-31")
	//                    }
	//            }
	//    },
	//    {
	//        $unwind: "$services"
	//    },
	//    {
	//        $group:
	//            {
	//                _id: { id: "$doctor._id", surname: "$doctor.surname", name: { $concat: [{ $substrCP: ["$doctor.name", 0, 1] }, "."] }, middlename: { $concat: [{ $substrCP: ["$doctor.middlename", 0, 1] }, "."] } },
	//                total: {
	//                    "$sum":
	//                        { "$toDouble": "$services.price" }
	//                }
	//            }
	//    }
	//]
	var pipeline = mongo.Pipeline{
		{{"$match", bson.D{{"dateEvent", bson.D{{"$gte", startDateEventB}, {"$lte", endDateEventB}}}}}},
		{{"$unwind", bson.D{{"path", "$services"}}}},
		{{"$group", bson.D{{"_id", bson.D{{"id", "$doctor._id"}, {"surname", "$doctor.surname"}, {"name", bson.A{{???}}}}}, {"total", bson.D{{"$sum", bson.D{{"$toDouble", "$services.price"}}}}}}}},
	}

	visitsCollection := s.MongoClient.Database("medregDB").Collection("visits")
	cursor, err := visitsCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var visitsBson []ProfitByDoctorModel
	if err = cursor.All(context.Background(), &visitsBson); err != nil {
		return nil, err
	}

	//caser := cases.Title(language.Russian)

	return res, nil
}

func (s *Mix) IncrementSeq(idx int, code string) error {
	sequenceCollection := s.MongoClient.Database("medregDB").Collection("sequence")
	filter := bson.M{"code": code}
	update := bson.D{{"$set", bson.D{{"seq", idx}}}}
	_, err := sequenceCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (s *Mix) GetSeq(code string) (int, error) {
	filter := bson.M{"code": code}
	sequenceCollection := s.MongoClient.Database("medregDB").Collection("sequence")

	var seq SequenceModel
	if err := sequenceCollection.FindOne(context.Background(), filter).Decode(&seq); err != nil {
		return -1, err
	}

	return seq.Seq, nil
}

func (s *Mix) Close() error {
	return s.MongoClient.Disconnect(context.Background())
}
