package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/booking", RequestBooking).Methods("POST")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func RequestBooking(w http.ResponseWriter, r *http.Request) {
	var bookingReq BookingReq
	_ = json.NewDecoder(r.Body).Decode(&bookingReq)
	code := StringWithCharset(5, charset)
	booking := Booking{code, bookingReq.Username, bookingReq.Destination}
	response := Response{"Booking Success, Your Booking Code :" + booking.Code}
	json.NewEncoder(w).Encode(response)
}

type BookingReq struct {
	Username    string `json:"username"`
	Destination string `json:"destination"`
}

type Booking struct {
	Code        string
	Username    string
	Destination string
}

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type Response struct {
	Message string
}
