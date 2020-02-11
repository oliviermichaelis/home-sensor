package infrastructure

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/oliviermichaelis/home-sensor/pkg/domain"
	"time"
)

type influxCloudHandler struct {
	client *influxdb2.Client
	org    string
	logger Logger
}

func NewInfluxCloudHandler(url string, token string, org string) (*influxCloudHandler, error) {
	influx := influxCloudHandler{
		logger: Logger{},
		org:    org,
	}

	client, err := influxdb2.New(url, token)
	if err != nil {
		message := fmt.Sprintf("influxdbcloud: %v", err)
		influx.logger.Fatal(message)
	}
	influx.client = client

	return &influx, nil
}

func (handler *influxCloudHandler) Insert(measurement domain.Measurement) {
	tags := map[string]string{
		"station": measurement.Station,
	}

	fields := map[string]interface{}{
		"temperature": measurement.Temperature,
		"humidity":    measurement.Humidity,
		"pressure":    measurement.Pressure,
	}

	timestamp, err := time.Parse("20060102150405", measurement.Timestamp)
	if err != nil {
		handler.logger.Log("influxdbcloud: could not parse timestamp: ", err)
		return
	}

	database, err := GetConfig(EnvInfluxDatabase)
	if err != nil {
		message := fmt.Sprintf("influxdbcloud: %v", err)
		handler.logger.Log(message)
	}

	data := []influxdb2.Metric{influxdb2.NewRowMetric(fields, database, tags, timestamp)}
	if _, err := handler.client.Write(context.Background(), "eosander", handler.org, data...); err != nil {
		message := fmt.Sprintf("influxdbcloud: %v", err)
		handler.logger.Log(message)
	}
}
