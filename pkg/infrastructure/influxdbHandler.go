package infrastructure

import (
	"fmt"
	influxdb "github.com/influxdata/influxdb1-client/v2"
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
	"log"
	"time"
)

// TODO refactor logging into struct
var logger = Logger{}
var measurementBatchPoint = influxdb.BatchPointsConfig {
	Precision:        "",
	Database:         "sensor",
	RetentionPolicy:  "",
	WriteConsistency: "",
}

type influxdbHandler struct {
	client influxdb.Client
}

func NewInfluxdbHandler(host string, port string, username string, password string) *influxdbHandler {
	// TODO implement TLS
	config := influxdb.HTTPConfig {
		Addr:               influxURL(host, port),
		Username:           username,
		Password:           password,
		UserAgent:          "",
		Timeout:            0,
		InsecureSkipVerify: false,
		TLSConfig:          nil,
		Proxy:              nil,
	}

	client , err := influxdb.NewHTTPClient(config)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return &influxdbHandler{client: client}
}

func (handler *influxdbHandler) Insert(measurement domain.Measurement) {
	// TODO fill in influxdb stuff. Recover from failures
	// TODO redial

	// transform measurement to influxdb point
	tags := map[string]string{
			"station": measurement.Station,
	}
	fields := map[string]interface{}{
		"temperature":	measurement.Temperature,
		"humidity":		measurement.Humidity,
		"pressure": 	measurement.Pressure,
	}
	timestamp, err := time.Parse("20060102150405", measurement.Timestamp)
	if err != nil {
		logger.Log("could not parse timestamp: ", err)
		return
	}

	points, err := influxdb.NewBatchPoints(measurementBatchPoint)
	if err != nil {
		logger.Fatal(err)
	}

	point, err := influxdb.NewPoint("sensor", tags, fields, timestamp)
	if err != nil {
		log.Fatal(err)
	}
	points.AddPoint(point)

	if err = handler.client.Write(points); err != nil {
		log.Fatal(err)
	}


	message := fmt.Sprintf("stored: %v", measurement)
	logger.Log(message)

	// TODO handle error in case unrecoverable

	// TODO in best case, password and username are re-read from file
}

// Returns the connection URL for the influxdb client
func influxURL(host string, port string) string {
	return fmt.Sprintf("http://%s:%s", host, port)
}
