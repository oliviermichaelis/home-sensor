package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"net/http/httptest"
	"testing"
)

func TestClimateHandlerInvalidMethod(t *testing.T) {
	request := httptest.NewRequest("GET", "https://apiserver.lab.oliviermichaelis.dev/measurements/climate", nil)
	recorder := httptest.NewRecorder()
	climateHandler(recorder, request)

	response := recorder.Result()
	fmt.Println(response.Body)
	if response.StatusCode != 405 {
		t.Errorf("Expected StatusCode 405, was: %v", response.StatusCode)
	}
}

func TestClimateHandlerNoBody(t *testing.T) {
	request := httptest.NewRequest("POST", "https://apiserver.lab.oliviermichaelis.dev/measurements/climate", nil)
	recorder := httptest.NewRecorder()
	climateHandler(recorder, request)

	response := recorder.Result()
	if response.StatusCode != 400 {
		t.Errorf("Expected StatusCode 400, was: %v", response.StatusCode)
	}
}

func TestClimateHandlerInvalidData(t *testing.T) {
	measurement := environment.SensorValues{
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
	climateHandler(recorder, request)

	response := recorder.Result()
	if response.StatusCode != 400 {
		t.Errorf("Expected StatusCode 400, was: %v", response.StatusCode)
	}
}

func TestClimateHandlerValidData(t *testing.T) {
	measurement := environment.SensorValues{
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
	climateHandler(recorder, request)

	response := recorder.Result()
	if response.StatusCode != 200 {
		t.Errorf("Expected StatusCode 400, was: %v", response.StatusCode)
	}
}
