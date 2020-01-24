package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
	"net/http/httptest"
	"testing"
)

type mockedMeasurementInteractor struct {}

func (m *mockedMeasurementInteractor) Store(measurement domain.Measurement) error {
	return nil
}

func TestClimateHandlerInvalidMethod(t *testing.T) {
	request := httptest.NewRequest("GET", "https://apiserver.lab.oliviermichaelis.dev/measurements/climate", nil)
	recorder := httptest.NewRecorder()

	webserviceHandler := WebserviceHandler{}
	webserviceHandler.ClimateHandler(recorder, request)

	response := recorder.Result()
	fmt.Println(response.Body)
	if response.StatusCode != 405 {
		t.Errorf("Expected StatusCode 405, was: %v", response.StatusCode)
	}
}

func TestClimateHandlerNoBody(t *testing.T) {
	request := httptest.NewRequest("POST", "https://apiserver.lab.oliviermichaelis.dev/measurements/climate", nil)
	recorder := httptest.NewRecorder()

	webserviceHandler := WebserviceHandler{}
	webserviceHandler.ClimateHandler(recorder, request)

	response := recorder.Result()
	if response.StatusCode != 400 {
		t.Errorf("Expected StatusCode 400, was: %v", response.StatusCode)
	}
}

func TestClimateHandlerInvalidData(t *testing.T) {
	measurement := domain.Measurement{
		Timestamp:   "",
		Station:     "",
		Temperature: 0,
		Humidity:    0,
		Pressure:    0,
	}

	body, _ := json.Marshal(measurement)
	reader := bytes.NewReader(body)
	request := httptest.NewRequest("POST", "https://apiserver.lab.oliviermichaelis.dev/measurements/climate", reader)
	recorder := httptest.NewRecorder()

	webserviceHandler := WebserviceHandler{}
	webserviceHandler.ClimateHandler(recorder, request)

	response := recorder.Result()
	if response.StatusCode != 400 {
		t.Errorf("Expected StatusCode 400, was: %v", response.StatusCode)
	}
}

func TestClimateHandlerValidData(t *testing.T) {
	measurement := domain.Measurement{
		Timestamp:   "20060102150405",
		Station:     "test",
		Temperature: 21.0,
		Humidity:    50.0,
		Pressure:    1024.5,
	}

	body, _ := json.Marshal(measurement)
	reader := bytes.NewReader(body)
	request := httptest.NewRequest("POST", "https://apiserver.lab.oliviermichaelis.dev/measurements/climate", reader)
	recorder := httptest.NewRecorder()

	webserviceHandler := WebserviceHandler{}
	webserviceHandler.MeasurementInteractor = &mockedMeasurementInteractor{}
	webserviceHandler.ClimateHandler(recorder, request)

	response := recorder.Result()
	if response.StatusCode != 200 {
		t.Errorf("Expected StatusCode 400, was: %v", response.StatusCode)
	}
}
