package domain

import (
	"errors"
	"time"
)

type MeasurementRepository interface {
	Store(measurement Measurement)
}

// Main struct that holds the measured values
type Measurement struct {
	Timestamp 	string	`json:"timestamp"`
	Station		string	`json:"station"`
	Temperature float64	`json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
}

func (s *Measurement) IsValid() error {
	if _, err := time.Parse("20060102150405", s.Timestamp); err != nil {
		return errors.New("Measurement: Timestamp is invalid")
	}

	if s.Station == "" {
		return errors.New("Measurement: Station is empty")
	}

	if -40.0 > s.Temperature || 100.0 < s.Temperature {
		return errors.New("Measurement: Temperature value is invalid")
	}

	if 0.0 > s.Humidity || s.Humidity > 100.0 {
		return errors.New("Measurement: Humidity value is invalid")
	}

	if s.Pressure < 900.0 || s.Pressure > 1100.0{
		return errors.New("Measurement: Pressure value is invalid")
	}

	return nil
}

