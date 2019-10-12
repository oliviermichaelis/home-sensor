package environment

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type SensorValues struct {
	Timestamp 	string
	Temperature float64
	Humidity    float64
	Pressure    float64
}

// Environment variables
var (
	rURL     = GetEnv("RABBITMQ_SERVICE_URL", "rabbitmq-ha.default.svc.cluster.local")
	rPort    = GetEnv("RABBITMQ_SERVICE_PORT", "5672")
	queue    = GetEnv("RABBITMQ_QUEUE", "sensor")
	exchange = GetEnv("RABBITMQ_EXCHANGE", "sensor")
	rSecret  = GetEnv("RABBITMQ_SECRET_PATH", "/passwords")
	iURL     = GetEnv("INFLUX_SERVICE_URL", "localhost")
	iPort    = GetEnv("INFLUX_SERVICE_PORT", "8086")
)

func RabbitmqURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", readUsername(rSecret), readPassword(rSecret), rURL, rPort)
}

// Returns the connection URL for the influxdb client
func InfluxURL() string {
	return fmt.Sprintf("http://%s:%s", iURL, iPort)
}

func readUsername(secretPath string) string {
	u, err := ioutil.ReadFile(secretPath + "/username")
	if err != nil {
		log.Fatal(err)
	}
	return string(u)
}

func readPassword(secretPath string) string {
	p, err := ioutil.ReadFile(secretPath + "/password")
	if err != nil {
		log.Fatal(err)
	}
	return string(p)
}

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
