package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
	"github.com/oliviermichaelis/home-sensor/pkg/infrastructure"
	"io/ioutil"
	"net/http"
)

type MeasurementInteractor interface {
	Store(measurement domain.Measurement) error
}

type WebserviceHandler struct {
	MeasurementInteractor MeasurementInteractor
	Logger                infrastructure.Logger
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
		message := fmt.Sprintf("Error to be handled is unknown: %d", status)
		handler.Logger.Fatal(message)
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
		handler.Logger.Log(fmt.Sprintf("Error reading body: %v", err))
		handler.errorHandler(w, r, http.StatusBadRequest)
		return
	}

	var measurement domain.Measurement
	if err := json.Unmarshal(body, &measurement); err != nil {
		handler.Logger.Log("webservice:", err)
		handler.errorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Persist measurement to whatever infrastructure implements the interface
	if err := handler.MeasurementInteractor.Store(measurement); err != nil {
		handler.Logger.Log("webservice:", err)
	}

	if _, err = w.Write(body); err != nil {
		handler.errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}
