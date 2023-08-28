package model

import "go.mongodb.org/mongo-driver/bson/primitive"

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
