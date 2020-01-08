package main

import (
	"encoding/json"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"io/ioutil"
	"log"
	"net/http"
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

	var measurement environment.SensorValues
	if err := json.Unmarshal(body, &measurement); err != nil {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	// validate the measurement struct
	if !measurement.IsValid() {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	// TODO use channel to send content to other goroutine which handles buffer.
	// Buffer could be queue of items with corresponding timeout. If item can't be delivered after 60min
	// the buffer controller should delete the message. If message delivery fails before timeout, the message
	// should be added to the tail.

}