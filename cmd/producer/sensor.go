package main

import (
	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"log"
	"math"
)

//TODO add timestamp
type SensorValues struct {
	Temperature float64
	Humidity    float64
	Pressure    float64
}

func setupSensor() {
	// Change loglevel of packages to omit debug output
	if err := logger.ChangePackageLogLevel("i2c", logger.WarnLevel); err != nil {
		log.Fatal(err)
	}
	if err := logger.ChangePackageLogLevel("bsbmp", logger.WarnLevel); err != nil {
		log.Fatal(err)
	}

	// Create new connection to i2c-bus on 1 line with address 0x76.
	connection, err := i2c.NewI2C(uint8(*i2cAddr), *i2cBus)
	if err != nil {
		log.Fatalf("Could not connect to i2c-bus! \n%v", err)
	}

	// Create new sensor connection to the bme280 sensor
	if sensor, err = bsbmp.NewBMP(bsbmp.BME280, connection); err != nil {
		log.Fatalf("Could not create new sensor! \n%v", err)
	}

	if err = sensor.IsValidCoefficients(); err != nil {
		log.Fatal(err)
	}
}

// Read temperature in celsius degree
func readTemperature() float64 {
	t, err := sensor.ReadTemperatureC(bsbmp.ACCURACY_HIGH)
	if err != nil {
		log.Fatalf("Could not read temperature!\n%v", err)
	}
	return round(float64(t))
}

// Read humidity in percent
func readHumidity() float64 {
	supported, h, err := sensor.ReadHumidityRH(bsbmp.ACCURACY_HIGH)
	if err != nil {
		log.Fatalf("Could not read temperature!\n%v", err)
	} else if !supported {
		log.Fatal("Humidity reading is not supported!")
	}
	return round(float64(h))
}

// Read Pressure in pascal
func readPressure() float64 {
	p, err := sensor.ReadPressurePa(bsbmp.ACCURACY_LOW)
	if err != nil {
		log.Fatalf("Could not read Pressure!\n%v", err)
	}
	return round(float64(p))
}

func readSensor() SensorValues {
	values := SensorValues{
		Temperature: readTemperature(),
		Humidity:    readHumidity(),
		Pressure:    readPressure(),
	}

	return values
}

// This function rounds floats to 2 decimal places
func round(value float64) float64 {
	return math.Floor(value*100) / 100
}
