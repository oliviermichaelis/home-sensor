package interfaces

import (
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
	"testing"
	"time"
)

type mockedDatabaseHandler struct{}

func (h *mockedDatabaseHandler) Insert(measurement domain.Measurement) {}

func (h *mockedDatabaseHandler) RetrieveLastWindow(station string, duration time.Duration) (*[]domain.Measurement, error) {
	return nil, nil
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
	databaseMeasurementRepo.influxCloudHandler = &mockedDatabaseHandler{}

	measurement := domain.Measurement{}
	measurement.PopulateTestValues()
	databaseMeasurementRepo.Store(measurement)
}
