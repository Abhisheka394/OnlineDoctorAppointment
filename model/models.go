package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId         string             `json:"userid" bson:"userid"`
	FirstName      string             `json:"firstname" bson:"firstname"`
	LastName       string             `json:"lastname" bson:"lastname"`
	Password       string             `json:"password" bson:"password"`
	UserType       string             `json:"usertype" bson:"usertype"`
	Age            int                `json:"age" bson:"age"`
	Specialization string             `json:"specialization" bson:"specialization"`
}

type Doctor struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId         string             `json:"userid" bson:"userid"`
	Name           string             `json:"firstname" bson:"firstname"`
	Age            int                `json:"age" bson:"age"`
	Specialization string             `json:"specialization" bson:"specialization"`
}

type Patient struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId string             `json:"userid" bson:"userid"`
	Name   string             `json:"name" bson:"name"`
	Age    int                `json:"age" bson:"age"`
}

type Appointment struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BookingId string             `json:"bookingid" bson:"bookingid"`
	DoctorId  string             `json:"doctorid" bson:"doctorid"`
	PatientId string             `json:"patientid" bson:"patientid"`
	StartTime int                `json:"starttime" bson:"starttime"`
	//EndTime   time.Time          `json:"endtime" bson:"endtime"`
	//Date    time.Time `json:"date" bson:"date"`
	Date string `json:"date" bson:"date"`
}

type Slot struct {
	StartTime int    `json:"starttime" bson:"starttime"`
	Date      string `json:"date" bson:"date"`
}
