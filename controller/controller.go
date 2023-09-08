package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"main.go/model"
)

//const connectionString = "mongodb+srv://abhi394394:Data1234@cluster0.2pfxcgr.mongodb.net/?retryWrites=true&w=majority"

//mongoURL= os.Getenv("MONGODB_URL")

const connectionString = "mongodb://localhost:27017"

var client *mongo.Client

func init() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
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
// http://localhost:9091/user/getdoc?specialization=<Write_Specialization_here>
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

func GetAvailability(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := client.Database("ODA").Collection("appointments")
	doctorid := r.URL.Query().Get("doctorid")
	date := r.URL.Query().Get("date")
	//starttime, _ := strconv.Atoi(r.URL.Query().Get("starttime"))
	var result []string

	for slottime := 9; slottime < 18; slottime++ {
		filter := bson.M{
			"doctorid":  doctorid,
			"date":      date,
			"starttime": slottime,
		}

		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			return
		}
		defer cursor.Close(ctx)
		var count = 0
		for cursor.Next(ctx) {
			count++
		}
		if err := cursor.Err(); err != nil {
			return
		}
		val := 4 - count
		var str string
		if slottime < 12 {
			str = "Time: " + fmt.Sprint(slottime) + "AM , No of slots available: " + fmt.Sprint(val)
		} else if slottime == 12 {
			str = "Time: " + fmt.Sprint(slottime) + "PM , No of slots available: " + fmt.Sprint(val)
		} else {
			str = "Time: " + fmt.Sprint(slottime-12) + "PM , No of slots available: " + fmt.Sprint(val)
		}
		result = append(result, str)

	}

	json.NewEncoder(w).Encode(result)

}

func CancelBooking(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := client.Database("ODA").Collection("appointments")
	doctorid := r.URL.Query().Get("doctorid")
	date := r.URL.Query().Get("date")
	patientid := r.URL.Query().Get("patientid")
	filter := bson.M{
		"doctorid":  doctorid,
		"date":      date,
		"patientid": patientid,
	}

	var dbAppointment model.Appointment
	err := collection.FindOne(ctx, filter).Decode(&dbAppointment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"Appointment doesn't exist"}`))
		return
	}
	cancelledAppointmentsCol := client.Database("ODA").Collection("cancelledAppointments")
	_, err = cancelledAppointmentsCol.InsertOne(ctx, dbAppointment)
	if err != nil {
		log.Fatal(err.Error())
	}

	deleteCount, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(deleteCount)
}
