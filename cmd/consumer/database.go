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

func setupClient() influxdb.Client {
	config := influxdb.HTTPConfig {
		Addr:               environment.InfluxURL(),	//TODO read from ENV
		Username:           "telegraf",			//TODO read from file
		Password:           "secretpassword",	//TODO read from file
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
	log.Printf("Connected successfully to %s", config.Addr)
	return client
}

func createDatabase(client influxdb.Client) {
	query := influxdb.NewQuery("CREATE DATABASE sensor", "", "")
	if response, err := client.Query(query); err != nil || response.Error() != nil {
		log.Fatalf("Could not create database! %v", err)
	}
}

func insertPoint(client influxdb.Client, values environment.SensorValues) {
	tags := map[string]string{}
	fields := map[string]interface{}{
		"temperature":	values.Temperature,
		"humidity":		values.Humidity,
		"pressure": 	values.Pressure,
	}
	timestamp, err := time.Parse("20060102150405", values.Timestamp)
	if err != nil {
		log.Fatalf("Could not parse timestamp! %v", err)
	}

	points, err := influxdb.NewBatchPoints(sensorBP)
	if err != nil {
		log.Fatal(err)
	}
	point, err := influxdb.NewPoint("sensor", tags, fields, timestamp)
	points.AddPoint(point)
	if err = client.Write(points);err != nil {
		log.Fatalf("Could not write to influxdb! %v", err)
	}
}