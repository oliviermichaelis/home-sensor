package domain

import (
	"testing"
)


func TestMeasurement_IsValid(t *testing.T) {
	measurement := Measurement{}
	measurement.PopulateTestValues()
	if err := measurement.IsValid(); err != nil {
		t.Error("unexpected error: ", err)
	}

	measurement.Timestamp = "200601021604050"
	if err := measurement.IsValid(); err == nil {
		t.Error("expected error, but is nil")
	}

	measurement = Measurement{}
	measurement.PopulateTestValues()
	measurement.Station = ""
	if err := measurement.IsValid(); err == nil {
		t.Error("expected error, but is nil")
	}

	measurement = Measurement{}
	measurement.PopulateTestValues()
	measurement.Temperature = -41.0
	if err := measurement.IsValid(); err == nil {
		t.Error("expected error, but is nil")
	}

	measurement = Measurement{}
	measurement.PopulateTestValues()
	measurement.Humidity = -1.0
	if err := measurement.IsValid(); err == nil {
		t.Error("expected error, but is nil")
	}

	measurement = Measurement{}
	measurement.PopulateTestValues()
	measurement.Humidity = 102.3
	if err := measurement.IsValid(); err == nil {
		t.Error("expected error, but is nil")
	}

	measurement = Measurement{}
	measurement.PopulateTestValues()
	measurement.Pressure = 1200.0
	if err := measurement.IsValid(); err == nil {
		t.Error("expected error, but is nil")
	}
}
