package main

import (
	"github.com/oliviermichaelis/home-sensor/pkg/infrastructure"
	"github.com/oliviermichaelis/home-sensor/pkg/interfaces"
	"github.com/oliviermichaelis/home-sensor/pkg/usecases"
	"net/http"
	"os"
)

func main() {
	// initialize logger
	logger := infrastructure.Logger{}

	// initialize configuration
	url, err := infrastructure.RegisterConfig("INFLUX_SERVICE_URL", "localhost")
	if err != nil {
		logger.Log(err.Error())
	}

	port, err := infrastructure.RegisterConfig("INFLUX_SERVICE_PORT", "8086")
	if err != nil {
		logger.Log(err.Error())
	}

	secretPath, err := infrastructure.RegisterConfig("INFLUX_SECRET_PATH", "/credentials/influx")
	if err != nil {
		logger.Log(err.Error())
	}

	if _, err := infrastructure.RegisterConfig("DEBUG", "false"); err != nil {
		logger.Log(err.Error())
	}

	if _, err := infrastructure.RegisterConfig("STATION_ID", ""); err != nil {
		logger.Log(err.Error())
	}

	username, err := infrastructure.ReadUsername(secretPath)
	if err != nil {
		logger.Log(err.Error())
		os.Exit(3)
	}

	password, err := infrastructure.ReadPassword(secretPath)
	if err != nil {
		logger.Log(err.Error())
		os.Exit(3)
	}

	// setup database connection
	databaseHandler := infrastructure.NewInfluxdbHandler(url, port, username, password)

	handlers := make(map[string]interfaces.DatabaseHandler)
	handlers["DatabaseMeasurementRepo"] = databaseHandler

	measurementInteractor := new(usecases.MeasurementInteractor)
	measurementInteractor.MeasurementRepository = interfaces.NewDatabaseMeasurementRepo(handlers)
	measurementInteractor.Logger = infrastructure.Logger{}

	webserviceHandler := interfaces.WebserviceHandler{
		MeasurementInteractor: measurementInteractor,
		Logger: infrastructure.Logger{},
	}

	http.HandleFunc("/measurements/climate", func(writer http.ResponseWriter, request *http.Request) {
		webserviceHandler.ClimateHandler(writer, request)
	})

	logger.Fatal(http.ListenAndServe(":8080", nil))
}
