package domain

import (
	"errors"
	"time"
)

type MeasurementRepository interface {
	Store(measurement Measurement)
}

// TODO change Timestamp to go time representation
// Main struct that holds the measured values
type Measurement struct {
	Timestamp 	string	`json:"timestamp"`
	Station		string	`json:"station"`
	Temperature float64	`json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
}

func (m *Measurement) IsValid() error {
	if _, err := time.Parse("20060102150405", m.Timestamp); err != nil {
		return errors.New("measurement: Timestamp is invalid")
	}

	if m.Station == "" {
		return errors.New("measurement: Station is empty")
	}

	if -40.0 > m.Temperature || 100.0 < m.Temperature {
		return errors.New("measurement: Temperature value is invalid")
	}

	if 0.0 > m.Humidity || m.Humidity > 100.0 {
		return errors.New("measurement: Humidity value is invalid")
	}

	if m.Pressure < 900.0 || m.Pressure > 1100.0{
		return errors.New("measurement: Pressure value is invalid")
	}

	return nil
}

func (m *Measurement) PopulateTestValues() {
	m.Timestamp = time.Now().UTC().Format("20060102150405")
	m.Station = "generated"
	m.Temperature = 20.28
	m.Humidity = 68.96
	m.Pressure = 1000.15
}


