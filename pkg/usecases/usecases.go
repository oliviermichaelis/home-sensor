package usecases

import (
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
	"time"
)

type Logger interface {
	Log(args ...interface{})
}

type MeasurementInteractor struct {
	MeasurementRepository domain.MeasurementRepository
	Logger                Logger
}

func (interactor *MeasurementInteractor) Store(measurement domain.Measurement) error {
	if err := measurement.IsValid(); err != nil {
		interactor.Logger.Log(err.Error())
		return err
	}

	interactor.MeasurementRepository.Store(measurement)
	return nil
}

func (interactor *MeasurementInteractor) RetrieveLastWindow(station string, duration time.Duration) (*[]domain.Measurement, error) {
	// TODO add input validation here, since this layer is tied to business logic. If input validation is made here, it's done for every outer layer
	return interactor.MeasurementRepository.RetrieveLastWindow(station, duration)
}
