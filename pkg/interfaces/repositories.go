package interfaces

import (
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
)

type DatabaseHandler interface {
	Insert(measurement domain.Measurement)
}

type DatabaseRepo struct {
	databaseHandlers 	map[string]DatabaseHandler
	databaseHandler		DatabaseHandler
}

type DatabaseMeasurementRepo DatabaseRepo

func NewDatabaseMeasurementRepo(databaseHandlers map[string]DatabaseHandler) *DatabaseMeasurementRepo {
	databaseMeasurementRepo := new(DatabaseMeasurementRepo)
	databaseMeasurementRepo.databaseHandlers = databaseHandlers
	databaseMeasurementRepo.databaseHandler = databaseHandlers["DatabaseMeasurementRepo"]
	return databaseMeasurementRepo
}

func (repo *DatabaseMeasurementRepo) Store(measurement domain.Measurement) {
	repo.databaseHandler.Insert(measurement)
}