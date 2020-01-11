package main

import (
	"encoding/json"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Handles error cases on API endpoints. Depending on the status code, the server responds with an explanation as json.
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	type Message struct {
		Subject string
		Message string
	}

	body := Message{}
	switch status {
	case http.StatusMethodNotAllowed:
		body = Message{
			Subject: "error",
			Message: "405 Method Not Allowed.",
		}
	case http.StatusBadRequest:
		body = Message{
			Subject: "error",
			Message: "400 Request is malformed",
		}
	default:
		log.Fatalf("Error to be handled is unknown: %d", status)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handles the /measurements/climate endpoint. Validates data and sends it to the goroutine handling the db insertion
func climateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Removes escaped double quotes from json generated in python
	bodyString, err := strconv.Unquote(string(body))
	if err == nil {
		body = []byte(bodyString)
	}

	var measurement environment.SensorValues
	if err := json.Unmarshal(body, &measurement); err != nil {
		log.Print(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Validate data in measurement struct
	if err := measurement.IsValid(); err != nil {
		log.Printf("Invalid Measurement from %s: %v", r.RemoteAddr, err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Send measurement to buffered channel which is consumed by readMeasurements()
	measurements <- measurement

	if _, err = w.Write(body); err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func setupServer() {
	http.HandleFunc("/measurements/climate", climateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
