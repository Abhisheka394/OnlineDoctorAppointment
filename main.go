package main

import (
	"fmt"
	"log"
	"net/http"

	"main.go/router"
)

func main() {
	fmt.Println("Sytem Started")
	r := router.Router()
	log.Fatal(http.ListenAndServe(":9090", r))
	fmt.Println("Listening at port 9090")
}
