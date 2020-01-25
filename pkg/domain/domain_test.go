package domain

import "testing"

var validMeasurement = Measurement{
	Timestamp:   "20060102160405",
	Station:     "test",
	Temperature: 20.28,
	Humidity:    58.95,
	Pressure:    1000.15,
}

func TestMeasurement_IsValid(t *testing.T) {
	measurement := validMeasurement
	if err := measurement.IsValid(); err != nil {
		t.Error("unexpected error: ", err)
	}

	measurement.Timestamp = "200601021604050"
	if err := measurement.IsValid(); err == nil {
		t.Error("expected error, but is nil")
	}

	measurement = validMeasurement
	measurement.Station = ""
	if err := measurement.IsValid(); err == nil {
		t.Error("expected error, but is nil")
	}

	measurement = validMeasurement
	measurement.Temperature = -41.0
	if err := measurement.IsValid(); err == nil {
		t.Error("expected error, but is nil")
	}

	measurement = validMeasurement
	measurement.Humidity = -1.0
	if err := measurement.IsValid(); err == nil {
		t.Error("expected error, but is nil")
	}

	measurement = validMeasurement
	measurement.Humidity = 102.3
	if err := measurement.IsValid(); err == nil {
		t.Error("expected error, but is nil")
	}

	measurement = validMeasurement
	measurement.Pressure = 1200.0
	if err := measurement.IsValid(); err == nil {
		t.Error("expected error, but is nil")
	}
}
