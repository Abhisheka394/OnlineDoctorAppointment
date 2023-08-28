package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"main.go/model"
)

const connectionString = "mongodb+srv://abhi394394:Data1234@cluster0.2pfxcgr.mongodb.net/?retryWrites=true&w=majority"

var SECRET_KEY = []byte("Victoria's Secret")

var client *mongo.Client

func init() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
}

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", err
	}
	return tokenString, nil
}

func UserSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	var dbUser model.User
	json.NewDecoder(r.Body).Decode(&user)
	user.UserId = strings.ToLower(user.UserId)
	user.Password = getHash([]byte(user.Password))
	collection := client.Database("ODA").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, bson.M{"userid": user.UserId}).Decode(&dbUser)
	if err == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"UserId already taken"}`))
		return
	}
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(w).Encode(result)

	if strings.ToLower(user.UserType) == "doctor" {
		var docUser model.Doctor
		docUser.UserId = user.UserId
		docUser.Name = user.FirstName + " " + user.LastName
		docUser.Specialization = user.Specialization
		docUser.Age = user.Age
		doctorCol := client.Database("ODA").Collection("doctor")
		res, _ := doctorCol.InsertOne(ctx, docUser)
		json.NewEncoder(w).Encode(res)
	} else {
		var patientUser model.Patient
		patientUser.UserId = user.UserId
		patientUser.Name = user.FirstName + " " + user.LastName
		patientUser.Age = user.Age
		patientCol := client.Database("ODA").Collection("patient")
		res, _ := patientCol.InsertOne(ctx, patientUser)
		json.NewEncoder(w).Encode(res)
	}

}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	var dbUser model.User
	json.NewDecoder(r.Body).Decode(&user)
	user.UserId = strings.ToLower(user.UserId)
	collection := client.Database("ODA").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, bson.M{"userid": user.UserId}).Decode(&dbUser)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}

	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)
	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil {
		log.Println(passErr)
		w.Write([]byte(`{"response":"Wrong Password!"}`))
		return
	}

	jwtToken, err := GenerateJWT()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}

	w.Write([]byte(`{"token":"` + jwtToken + `"}`))
}
