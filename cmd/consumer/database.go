package main

import (
	influxdb "github.com/influxdata/influxdb1-client/v2"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"log"
	"time"
)

var sensorBP = influxdb.BatchPointsConfig {
	Precision:        "",
	Database:         "sensor",
	RetentionPolicy:  "",
	WriteConsistency: "",
}

// Sets up a infludcb client
func setupClient() influxdb.Client {
	username, err := environment.ReadUsername(environment.ISecret)
	if err != nil {
		log.Fatalf("Could not read username for influxdb! %v", err)
	}
	password, err := environment.ReadPassword(environment.ISecret)
	if err != nil {
		log.Fatalf("Could not read password for influxdb! %v", err)
	}
	config := influxdb.HTTPConfig {
		Addr:               environment.InfluxURL(),
		Username:           username,
		Password:           password,
		UserAgent:          "",
		Timeout:            0,
		InsecureSkipVerify: false,
		TLSConfig:          nil,
		Proxy:              nil,
	}

	client, err := influxdb.NewHTTPClient(config)
	if err != nil {
		log.Fatalf("Could not instantiate influxdb client! %v", err)
	}
	log.Printf("Successfully setup influxdb client to %s", config.Addr)
	return client
}

// This function tries to create a database in case it doesn't exist yet.
// That's also the first time the client actually connects to the database.
func createDatabase(client influxdb.Client) {
	query := influxdb.NewQuery("CREATE DATABASE sensor", "", "")
	if response, err := client.Query(query); err != nil || response.Error() != nil {
		log.Fatalf("Could not create database! %v", err)
	}
}

// Inserts a single timeseries point into the database
func insertPoint(client influxdb.Client, values environment.SensorValues) error {
	tags := map[string]string{
		"station": values.Station,
	}
	fields := map[string]interface{}{
		"temperature":	values.Temperature,
		"humidity":		values.Humidity,
		"pressure": 	values.Pressure,
	}
	timestamp, err := time.Parse("20060102150405", values.Timestamp)
	if err != nil {
		return err
	}

	points, err := influxdb.NewBatchPoints(sensorBP)
	if err != nil {
		return err
	}
	point, err := influxdb.NewPoint("sensor", tags, fields, timestamp)
	points.AddPoint(point)
	if err = client.Write(points);err != nil {
		return err
	}

	return nil
}