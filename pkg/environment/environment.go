package environment

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Environment variables
var (
	serviceURL = GetEnv("RABBITMQ_SERVICE_URL", "rabbitmq-ha.default.svc.cluster.local")
	servicePort = GetEnv("RABBITMQ_SERVICE_PORT", "5672")
	queue = GetEnv("RABBITMQ_QUEUE", "sensor")
	exchange = GetEnv("RABBITMQ_EXCHANGE", "sensor")
	secretPath = GetEnv("SECRET_PATH", "/passwords")
)

func AssembleURL() string {
	username, err := ioutil.ReadFile(secretPath + "/username")
	if err != nil {
		log.Fatal(err)
	}
	password, err := ioutil.ReadFile(secretPath + "/password")
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", string(username), string(password), serviceURL, servicePort)
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
