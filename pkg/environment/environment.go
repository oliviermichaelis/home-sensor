package environment

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
	rSecret  = GetEnv("RABBITMQ_SECRET_PATH", "/credentials/rabbitmq")
	iURL     = GetEnv("INFLUX_SERVICE_URL", "localhost")
	iPort    = GetEnv("INFLUX_SERVICE_PORT", "8086")
	ISecret  = GetEnv("INFLUX_SECRET_PATH", "/credentials/influx")
	Debug,_	 = strconv.ParseBool(GetEnv("DEBUG", "false"))
)

func RabbitmqURL() string {
	log.Printf("rSecret: %s", rSecret)
	username, err := ReadUsername(rSecret)
	log.Printf("user: %s", username)
	if err != nil {
		log.Fatalf("Couldn't read username: %v", err)
	}

	password, err := ReadUsername(rSecret)
	log.Printf("pass: %s", password)
	if err != nil {
		log.Fatalf("Couldn't read password: %v", err)
	}

	return fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, rURL, rPort)
}

// Returns the connection URL for the influxdb client
func InfluxURL() string {
	return fmt.Sprintf("http://%s:%s", iURL, iPort)
}

func ReadUsername(secretPath string) (string, error) {
	u, err := ioutil.ReadFile(secretPath + "/username")
	return string(u), err
}

func ReadPassword(secretPath string) (string, error) {
	log.Println("Testing!")
	t1, err := ioutil.ReadFile(secretPath + "/password")
	log.Printf("t1: %s", t1)
	t2, err := ioutil.ReadFile("/credentials/rabbitmq/password")
	log.Printf("t2: %s", t2)
	path := secretPath + "/password"
	log.Println("t3 path: ", path)
	t3, err := ioutil.ReadFile(path)
	log.Printf("t3: %s", t3)

	p, err := ioutil.ReadFile(secretPath + "/password")
	return string(p), err
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
