package store

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

type VisitModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Client    ClientModel        `bson:"client,omitempty"`
	DateEvent primitive.DateTime `bson:"dateEvent,omitempty"`
	Services  []ServiceModel     `bson:"services,omitempty"`
	TotalSum  int
}

type ClientModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Surname    string             `bson:"surname,omitempty"`
	Middlename string             `bson:"middlename,omitempty"`
	Gender     string             `bson:"gender,omitempty"`
}

type ServiceModel struct {
	ID       string `bson:"_id,omitempty"`
	Price    string `bson:"price,omitempty"`
	Discount string `bson:"discount,omitempty"`
	Category string `bson:"category,omitempty"`
}

type SequenceModel struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Code string             `bson:"code,omitempty"`
	Seq  int                `bson:"seq,omitempty"`
}

func (v *VisitModel) CalculateTotalSum() int {
	res := 0
	for _, s := range v.Services {
		p, err := strconv.Atoi(s.Price)
		if err != nil {
			panic(err)
		}
		d, err := strconv.Atoi(s.Discount)
		if err != nil {
			panic(err)
		}
		x := 1.0

		if d > 0 {
			x = (100.0 - float64(d)) / 100.0
		}

		//Truncate penny in case there is discount * price not equal int number
		res = res + int(float64(p)*x)
	}

	return res
}
