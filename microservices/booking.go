package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

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
	SendMessage(booking)
	response := Response{"Booking Success, Your Booking Code :" + booking.Code}
	json.NewEncoder(w).Encode(response)
}

func SendMessage(booking Booking) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect AMQP")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to connect Channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"booking",
		true,
		false, //no-
		false, //
		false, //
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msg, err := json.Marshal(booking)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(msg))
	err = ch.Publish(
		"notifExchange",
		q.Name,
		false, //
		false, //
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(string(msg)),
		})

	log.Printf("[x] Sent %s", msg)
	failOnError(err, "Failed to publish a messages")
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
