package usecases

import (
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
	"github.com/oliviermichaelis/home-sensor/pkg/infrastructure"
	"testing"
)

var validMeasurement = domain.Measurement{
	Timestamp:   "20060102160405",
	Station:     "test",
	Temperature: 20.28,
	Humidity:    58.95,
	Pressure:    1000.15,
}

type mockedMeasurementRepository struct{}

func (m *mockedMeasurementRepository) Store(measurement domain.Measurement) {
	return
}
func TestMeasurementInteractor_Store(t *testing.T) {
	measurementInteractor := MeasurementInteractor{
		MeasurementRepository: &mockedMeasurementRepository{},
		Logger:                infrastructure.Logger{},
	}

	if err := measurementInteractor.Store(validMeasurement); err != nil {
		t.Error("didnt expect error:", err)
	}

	invalidMeasurement := validMeasurement
	invalidMeasurement.Humidity = 123.4
	if err := measurementInteractor.Store(invalidMeasurement); err == nil {
		t.Error("expected error, not nil")
	}
}
