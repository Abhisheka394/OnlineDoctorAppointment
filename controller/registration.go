package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"main.go/model"
)

var strCh = make(chan string)

// Converting the string password to hash
func getHash(pwd []byte) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)

	}
	strCh <- string(hash)
	//return string(hash)
}

// For adding doctors to the doctor table during registration
func AddDoctor(user model.User, w http.ResponseWriter) {
	var docUser model.Doctor
	docUser.UserId = user.UserId
	docUser.Name = user.FirstName + " " + user.LastName
	docUser.Specialization = user.Specialization
	docUser.Age = user.Age
	doctorCol := client.Database("ODA").Collection("doctor")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, _ := doctorCol.InsertOne(ctx, docUser)
	json.NewEncoder(w).Encode(res)
}

// For adding patients to the patient table during registration
func AddPatient(user model.User, w http.ResponseWriter) {
	var patientUser model.Patient
	patientUser.UserId = user.UserId
	patientUser.Name = user.FirstName + " " + user.LastName
	patientUser.Age = user.Age
	patientCol := client.Database("ODA").Collection("patient")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, _ := patientCol.InsertOne(ctx, patientUser)
	json.NewEncoder(w).Encode(res)
}

// Validating the data fields
func Validate(user model.User, w http.ResponseWriter) bool {
	if user.UserId == "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"User Id cannot be blank"}`))
		return false
	} else if user.FirstName == "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"First Name cannot be blank"}`))
		return false
	} else if user.Age <= 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"Age cannot be less than 0"}`))
		return false
	} else if len(user.Password) < 5 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"Password length should be greater than 5"}`))
		return false
	} else if strings.ToLower(user.UserType) == "doctor" {
		if user.Specialization == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"Specialization cannot be empty for a doctor"}`))
			return false
		}
	}
	return true
}

// For Registring new users
func UserSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	//var dbUser model.User
	json.NewDecoder(r.Body).Decode(&user)

	if !Validate(user, w) {
		return
	}

	//go getHash([]byte(user.Password))
	user.UserId = strings.ToLower(user.UserId)
	go getHash([]byte(user.Password))
	collection := client.Database("ODA").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//err := collection.FindOne(ctx, bson.M{"userid": user.UserId})
	count, _ := collection.CountDocuments(context.TODO(), bson.M{"userid": user.UserId})
	if count != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"UserId already taken"}`))
		return
	}
	user.Password = <-strCh
	result, e := collection.InsertOne(ctx, user)
	if e != nil {
		fmt.Println(e)
	}
	//fmt.Println(result)
	json.NewEncoder(w).Encode(result)

	if strings.ToLower(user.UserType) == "doctor" {
		AddDoctor(user, w)
	} else {
		AddPatient(user, w)
	}

}
