package main

import (
	"errors"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"log"
)

// Global channel used as buffer
var measurements = make(chan environment.SensorValues, 1024)

func main() {
	go setupServer()
	log.Fatal(readMeasurements(measurements))
}

func readMeasurements(messages <-chan environment.SensorValues) error {
	client := setupClient()
	defer client.Close()

	for message := range messages {
		//TODO if messages buffers more than 1 message, aggregate and send multiple points at once
		if err := insertPoint(client, message); err != nil {
			return err
		}
		log.Printf("Inserted into influxdb: %v", message)
	}

	return errors.New("readMeasurements: infinite loop got broken")
}
