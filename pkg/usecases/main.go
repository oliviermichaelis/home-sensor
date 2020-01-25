package usecases

import (
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
)

type Logger interface {
	Log(args ...interface{})
}

type MeasurementInteractor struct {
	MeasurementRepository domain.MeasurementRepository
	Logger Logger
}

func (interactor *MeasurementInteractor) Store(measurement domain.Measurement) error {
	if err := measurement.IsValid(); err != nil {
		interactor.Logger.Log(err.Error())
		return err
	}

	interactor.MeasurementRepository.Store(measurement)
	return nil
}
