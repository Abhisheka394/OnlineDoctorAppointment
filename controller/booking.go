package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"main.go/model"
)

func validateBookingDetails(doctorid string, date string, starttime int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("ODA").Collection("appointments")
	var result []model.Appointment

	filter := bson.M{
		"doctorid":  doctorid,
		"date":      date,
		"starttime": starttime,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var temp model.Appointment
		if err := cursor.Decode(&temp); err != nil {
			return err
		}
		result = append(result, temp)
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	if len(result) >= 4 {
		return errors.New("the selected slot is full")
	}

	date_format := "2006-01-02"
	parsed_date, _ := time.Parse(date_format, date)
	for i := 0; i < starttime-6; i++ {
		parsed_date = parsed_date.Add(time.Hour)
	}
	parsed_date = parsed_date.Add(30 * time.Minute)

	// fmt.Println(parsed_date)
	// fmt.Println(time.Now())
	if parsed_date.Before(time.Now()) {
		return errors.New("the Booking time for the entered date has already passed")
	}

	return nil

}

// For Booking appointments
// url : http://localhost:9091/user/booking?doctorid=<>&patientid=<>
// json : { "starttime":, "date":""}
func BookAppointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var booking model.Appointment
	var slot model.Slot
	booking.DoctorId = r.URL.Query().Get("doctorid")
	booking.PatientId = r.URL.Query().Get("patientid")
	json.NewDecoder(r.Body).Decode(&slot)
	booking.StartTime = slot.StartTime
	booking.Date = slot.Date
	// date_format := "2006-01-02"
	// parsed_date, _ := time.Parse(date_format, slot.Date)

	booking.BookingId = booking.PatientId + booking.Date

	err := validateBookingDetails(booking.DoctorId, booking.Date, booking.StartTime)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("ODA").Collection("appointments")
	result, err := collection.InsertOne(ctx, booking)
	if err != nil {
		log.Fatal(err.Error())
	}
	json.NewEncoder(w).Encode(result)

}
