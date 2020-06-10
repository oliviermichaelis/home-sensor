package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type ClimateEndpoint struct {
	Repositories []Repository
}

// Handles the /measurements/climate endpoint. Validates data and sends it to the goroutine handling the db connection
func (c *ClimateEndpoint) climateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c.climateGetHandler(w, r)
	case "POST":
		c.climatePostHandler(w, r)
	default:
		errorHandler(w, r, http.StatusMethodNotAllowed)
	}
}

/*
This endpoints accepts a GET request to /measurement/climate with the following parameters:
station: the name of the station
duration: window size in seconds
*/
func (c *ClimateEndpoint) climateGetHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	var station, d string
	if ok := parseQueryParameter(w, r, params, "station", &station); !ok {
		return
	}
	if ok := parseQueryParameter(w, r, params, "duration", &d); !ok {
		return
	}

	duration, err := time.ParseDuration(d + "s")
	if err != nil {
		fmt.Printf("climategethandler: %v", err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	now := time.Now()
	window, err := c.Repositories[0].retrieveWindow(station, now.Add(-duration), now)
	if err != nil {
		fmt.Printf("climategethandler: %v", err)
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	// if window is nil, we need to inform client
	if window == nil {
		fmt.Println("climategethandler: returned window is nil")
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	// TODO refactor into reusable function
	// send data as response to client
	j, err := json.Marshal(window)
	if err != nil {
		fmt.Printf("climategethandler: %v", err)
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(j); err != nil {
		fmt.Printf("climategethandler: %v", err)
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func (c *ClimateEndpoint) climatePostHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	var measurement measurement
	if err := json.Unmarshal(body, &measurement); err != nil {
		fmt.Printf("webservice: %v", err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Persist measurement to whatever infrastructure implements the interface
	if err := c.Repositories[0].insert(measurement); err != nil {
		fmt.Printf("webservice: %v", err)
	}

	if _, err = w.Write(body); err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

// Handles error cases on API endpoints. Depending on the status code, the server responds with an explanation as json.
func errorHandler(w http.ResponseWriter, _ *http.Request, status int) {
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
		fmt.Println(message)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parseQueryParameter(w http.ResponseWriter, r *http.Request, params url.Values, key string, value *string) bool {
	v, ok := params[key]
	if !ok {
		fmt.Printf("/measurement/climate: %s query parameter is missing", key)
		errorHandler(w, r, http.StatusBadRequest)
	}
	*value = v[0]

	return true
}
