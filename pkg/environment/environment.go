package environment

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

// Main struct that holds the measured values
type SensorValues struct {
	Timestamp 	string	`json:"timestamp"`
	Station		string	`json:"station"`
	Temperature float64	`json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
}

// Environment variables
var (
	rURL     = GetEnv("RABBITMQ_SERVICE_URL", "rabbitmq-ha.default.svc.cluster.local")
	rPort    = GetEnv("RABBITMQ_SERVICE_PORT", "5672")
	queue    = GetEnv("RABBITMQ_QUEUE", "sensor")
	exchange = GetEnv("RABBITMQ_EXCHANGE", "sensor")
	rSecret  = GetEnv("RABBITMQ_SECRET_PATH", "/credentials/rabbitmq")
	iURL     = GetEnv("INFLUX_SERVICE_URL", "localhost")
	iPort    = GetEnv("INFLUX_SERVICE_PORT", "8086")
	ISecret  = GetEnv("INFLUX_SECRET_PATH", "/credentials/influx")
	Debug,_	 = strconv.ParseBool(GetEnv("DEBUG", "false"))
	Station  = GetEnv("STATION_ID", "")
)

func (s SensorValues) IsValid() error {
	if s == (SensorValues{}) {
		return errors.New("SensorValues: is initialized empty")
	}

	if _, err := time.Parse("20060102150405", s.Timestamp); err != nil {
		return errors.New("SensorValues: Timestamp is invalid")
	}

	if s.Station == "" {
		return errors.New("SensorValues: Station is empty")
	}

	if -40.0 > s.Temperature || 100.0 < s.Temperature {
		return errors.New("SensorValues: Temperature value is invalid")
	}

	if -0.1 > s.Humidity || s.Humidity > 100.0 {
		return errors.New("SensorValues: Humidity value is invalid")
	}

	if s.Pressure < 900.0 || s.Pressure > 1100.0{
		return errors.New("SensorValues: Pressure value is invalid")
	}

	return nil
}

func RabbitmqURL() string {
	username, err := ReadUsername(rSecret)
	if err != nil {
		log.Fatalf("Couldn't read username: %v", err)
	}

	password, err := ReadPassword(rSecret)
	if err != nil {
		log.Fatalf("Couldn't read password: %v", err)
	}

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, rURL, rPort)
}

// Returns the connection URL for the influxdb client
func InfluxURL() string {
	return fmt.Sprintf("http://%s:%s", iURL, iPort)
}

// Returns the username read from a 'username' file under the specified parameter
func ReadUsername(secretPath string) (string, error) {
	u, err := ioutil.ReadFile(secretPath + "/username")
	return string(u), err
}

// Returns the password read from a 'password' file under the specified parameter
func ReadPassword(secretPath string) (string, error) {
	p, err := ioutil.ReadFile(secretPath + "/password")
	return string(p), err
}

// Returns an environment variable if it exists, else it returns a fallback variable passed as parameter
func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func GetQueue() string {
	return queue
}

func GetExchange() string{
	return exchange
}
