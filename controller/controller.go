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

// Converting the string password to hash
func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// Generating JWT token to return after every successful login
func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", err
	}
	return tokenString, nil
}

// For adding doctors to the doctor table during registration
func AddDoctor(user model.User, w http.ResponseWriter) {
	var docUser model.Doctor
	docUser.UserId = user.UserId
	docUser.Name = user.FirstName + " " + user.LastName
	docUser.Specialization = user.Specialization
	docUser.Age = user.Age
	doctorCol := client.Database("ODA").Collection("doctor")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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

	user.UserId = strings.ToLower(user.UserId)
	user.Password = getHash([]byte(user.Password))
	collection := client.Database("ODA").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, bson.M{"userid": user.UserId})
	if err == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"UserId already taken"}`))
		return
	}
	result, _ := collection.InsertOne(ctx, user)
	//fmt.Println(result)
	json.NewEncoder(w).Encode(result)

	if strings.ToLower(user.UserType) == "doctor" {
		AddDoctor(user, w)
	} else {
		AddPatient(user, w)
	}

}

// For user Login
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

// Geting the details of doctors based on specialization
func getDoctorHelper(specialization string) ([]model.Doctor, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	doctorCol := client.Database("ODA").Collection("doctor")

	// Filter for query
	filter := bson.M{"specialization": specialization}

	cursor, err := doctorCol.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var doctors []model.Doctor
	for cursor.Next(ctx) {
		var doctor model.Doctor
		if err := cursor.Decode(&doctor); err != nil {
			return nil, err
		}
		doctors = append(doctors, doctor)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return doctors, nil

}

// Function to get doctor details based on specialization
// http://localhost:9090/user/getdoc?specialization=<Write_Specialization_here>
func GetDoctorDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	specialization := r.URL.Query().Get("specialization")
	doctors, err := getDoctorHelper(specialization)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"Error fetching doctors list"}`))
		return
	}
	json.NewEncoder(w).Encode(doctors)

}
