package interfaces

import (
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
	"testing"
)

type mockedDatabaseHandler struct {}

func (h *mockedDatabaseHandler) Insert(measurement domain.Measurement) {
	return
}

func TestNewDatabaseMeasurementRepo(t *testing.T) {
	databaseHandlers := make(map[string]DatabaseHandler)
	databaseMeasurementRepo := NewDatabaseMeasurementRepo(databaseHandlers)

	if databaseMeasurementRepo == nil {
		t.Error("repositories: new databasemeasurementrepo is nil")
	}
}

func TestDatabaseMeasurementRepo_Store(t *testing.T) {
	databaseHandlers := make(map[string]DatabaseHandler)
	databaseHandlers["DatabaseMeasurementRepo"] = &mockedDatabaseHandler{}
	databaseMeasurementRepo := NewDatabaseMeasurementRepo(databaseHandlers)
	databaseMeasurementRepo.databaseHandler = &mockedDatabaseHandler{}

	measurement := domain.Measurement{}
	measurement.PopulateRandomValues()
	databaseMeasurementRepo.Store(measurement)
}