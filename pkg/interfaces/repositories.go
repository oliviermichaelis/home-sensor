package interfaces

import (
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
	"time"
)

type DatabaseHandler interface {
	Insert(measurement domain.Measurement)
	RetrieveLastWindow(station string, duration time.Duration) (*[]domain.Measurement, error)
}

type DatabaseRepo struct {
	databaseHandlers   map[string]DatabaseHandler
	databaseHandler    DatabaseHandler
	influxCloudHandler DatabaseHandler
}

type DatabaseMeasurementRepo DatabaseRepo

func NewDatabaseMeasurementRepo(databaseHandlers map[string]DatabaseHandler) *DatabaseMeasurementRepo {
	databaseMeasurementRepo := new(DatabaseMeasurementRepo)
	databaseMeasurementRepo.databaseHandlers = databaseHandlers
	databaseMeasurementRepo.databaseHandler = databaseHandlers["DatabaseMeasurementRepo"]
	databaseMeasurementRepo.influxCloudHandler = databaseHandlers["DatabaseInfluxCloudRepo"]
	return databaseMeasurementRepo
}

func (repo *DatabaseMeasurementRepo) Store(measurement domain.Measurement) {
	repo.databaseHandler.Insert(measurement)
	repo.influxCloudHandler.Insert(measurement)
}

func (repo *DatabaseMeasurementRepo) RetrieveLastWindow(station string, duration time.Duration) (*[]domain.Measurement, error) {
	return repo.databaseHandler.RetrieveLastWindow(station, duration)
}
