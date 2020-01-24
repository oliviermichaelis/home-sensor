package interfaces

import (
	"encoding/json"
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type MeasurementInteractor interface {
	Store(measurement domain.Measurement) error
}

type WebserviceHandler struct {
	MeasurementInteractor MeasurementInteractor
}


// Handles error cases on API endpoints. Depending on the status code, the server responds with an explanation as json.
func (handler WebserviceHandler) errorHandler(w http.ResponseWriter, r *http.Request, status int) {
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
func (handler WebserviceHandler) ClimateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		handler.errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		handler.errorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Removes escaped double quotes from json generated in python
	bodyString, err := strconv.Unquote(string(body))
	if err == nil {
		body = []byte(bodyString)
	}

	var measurement domain.Measurement
	if err := json.Unmarshal(body, &measurement); err != nil {
		log.Print(err)
		handler.errorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Validate data in measurement struct
	if err := measurement.IsValid(); err != nil {
		log.Printf("Invalid Measurement from %s: %v", r.RemoteAddr, err)
		handler.errorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Persist measurement to whatever infrastructure implements the interface
	if err := handler.MeasurementInteractor.Store(measurement); err != nil {
		log.Println(err.Error())
	}

	if _, err = w.Write(body); err != nil {
		handler.errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}
