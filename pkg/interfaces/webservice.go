package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
	"github.com/oliviermichaelis/home-sensor/pkg/infrastructure"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type MeasurementInteractor interface {
	Store(measurement domain.Measurement) error
	RetrieveLastWindow(station string, duration time.Duration) (*[]domain.Measurement, error)
}

type WebserviceHandler struct {
	MeasurementInteractor MeasurementInteractor
	Logger                infrastructure.Logger
}

// Handles error cases on API endpoints. Depending on the status code, the server responds with an explanation as json.
func (handler WebserviceHandler) errorHandler(w http.ResponseWriter, _ *http.Request, status int) {
	type Message struct {
		Subject string
		Message string
	}

	body := Message{}
	switch status {
	case http.StatusBadRequest:
		body = Message{
			Subject: "error",
			Message: "400 Request is malformed",
		}
	case http.StatusNotFound:
		body = Message{
			Subject: "error",
			Message: "404 Not Found",
		}
	case http.StatusMethodNotAllowed:
		body = Message{
			Subject: "error",
			Message: "405 Method Not Allowed.",
		}
	case http.StatusInternalServerError:
		body = Message{
			Subject: "error",
			Message: "500 Internal Server Error",
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

// Handles the /measurements/climate endpoint. Validates data and sends it to the goroutine handling the db connection
func (handler WebserviceHandler) ClimateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handler.climateGetHandler(w, r)
	case "POST":
		handler.climatePostHandler(w, r)
	default:
		handler.errorHandler(w, r, http.StatusMethodNotAllowed)
	}
}

/*
This endpoints accepts a GET request to /measurement/climate with the following parameters:
station: the name of the station
duration: window size in seconds
 */
func (handler WebserviceHandler) climateGetHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	var station, d string
	if ok := handler.parseQueryParameter(w, r, params, "station", &station); !ok {
		return
	}
	if ok := handler.parseQueryParameter(w, r, params, "duration", &d); !ok {
		return
	}
	//if ok := handler.parseQueryParameter(w, r, params, "timestamp", &t); !ok {
	//	return
	//}

	duration, err := time.ParseDuration(d + "s")
	if err != nil {
		handler.Logger.Log("climategethandler:", err)
		handler.errorHandler(w, r, http.StatusBadRequest)
		return
	}

	//timestamp, err := time.Parse(time.RFC3339, t)
	//if err != nil {
	//	handler.Logger.Log("climategethandler:", err)
	//	handler.errorHandler(w, r, http.StatusBadRequest)
	//	return
	//}

	window, err := handler.MeasurementInteractor.RetrieveLastWindow(station, duration)
	if err != nil {
		handler.Logger.Log(fmt.Sprintf("climategethandler: %v", err))
		handler.errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	// if window is nil, we need to inform client
	if window == nil {
		handler.Logger.Log("climategethandler: returned window is nil")
		handler.errorHandler(w, r, http.StatusNotFound)
		return
	}

	// TODO refactor into reusable function
	// send data as response to client
	j, err := json.Marshal(window)
	if err != nil {
		handler.Logger.Log(fmt.Sprintf("climategethandler: %v", err))
		handler.errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(j); err != nil {
		handler.Logger.Log(fmt.Sprintf("climategethandler: %v", err))
		handler.errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func (handler WebserviceHandler) climatePostHandler(w http.ResponseWriter, r *http.Request) {
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

func (handler WebserviceHandler) parseQueryParameter(w http.ResponseWriter, r *http.Request, params url.Values, key string, value *string) bool {
	v, ok := params[key]
	if !ok {
		handler.Logger.Log(fmt.Sprintf("/measurement/climate: %s query parameter is missing", key))
		handler.errorHandler(w, r, http.StatusBadRequest)
	}
	*value = v[0]

	return true
}
