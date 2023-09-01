package router

import (
	"github.com/gorilla/mux"
	"main.go/controller"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/user/login", controller.UserLogin).Methods("POST")
	router.HandleFunc("/user/signUp", controller.UserSignup).Methods("POST")
	router.HandleFunc("/user/getdoc", controller.GetDoctorDetails).Methods("GET")
	router.HandleFunc("/user/booking", controller.BookAppointment).Methods("POST")
	router.HandleFunc("/user/getavailability", controller.GetAvailability).Methods("GET")
	router.HandleFunc("/user/cancelbooking", controller.CancelBooking).Methods("DELETE")
	return router
}
