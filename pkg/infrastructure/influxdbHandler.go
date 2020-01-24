package infrastructure

import (
	"fmt"
	influxdb "github.com/influxdata/influxdb1-client/v2"
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
	"net"
	"time"
)

var measurementBatchPoint = influxdb.BatchPointsConfig {
	Precision:        "",
	Database:         "sensor",
	RetentionPolicy:  "",
	WriteConsistency: "",
}
type influxdbHandler struct {
	client influxdb.Client
	logger Logger
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
	
	influx := influxdbHandler{
		logger: Logger{},
	}

	client , err := influxdb.NewHTTPClient(config)
	if err != nil {
		influx.logger.Fatal(err.Error())
	}
	influx.client = client
	return &influx
}

func (handler *influxdbHandler) Insert(measurement domain.Measurement) {
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
		handler.logger.Log("influxdb: could not parse timestamp: ", err)
		return
	}

	points, err := influxdb.NewBatchPoints(measurementBatchPoint)
	if err != nil {
		message := "influxdb: could not create batchpoints: "
		handler.logger.Log(message, err)
		return
	}

	point, err := influxdb.NewPoint("sensor", tags, fields, timestamp)
	if err != nil {
		message := "influxdb: could not create batchpoint: "
		handler.logger.Log(message, err)
	}
	points.AddPoint(point)

	// Each goroutine will try to reconnect every 10s in case the data could'nt be persisted to influxdb
	for {
		err = nil
		if err = handler.client.Write(points); err != nil {
			handler.logger.Log("influxdb:", err)
			if _, ok := err.(net.Error); !ok {
				// In case the returned error is not a net.Error, the error isn't a networking malfunction/error and can
				// therefore not be resolved by continuous retries
				handler.logger.Log("influxdb: error is not net.Error and retries are not prohibited")
				return
			}
		}

		// In case no error has occured, break out of infinite retry loop
		if err == nil {
			break
		}
		time.Sleep(time.Second * 10)
	}

	message := fmt.Sprintf("influxdb: successful insert: %v", measurement)
	handler.logger.Log(message)
}

// Returns the connection URL for the influxdb client
func influxURL(host string, port string) string {
	return fmt.Sprintf("http://%s:%s", host, port)
}
