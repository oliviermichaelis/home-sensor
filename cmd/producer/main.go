package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/d2r2/go-i2c"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// cli flags
var (
	samplerate = flag.Duration("samplerate", 60 * time.Second, "Sample rate at which sensor is polled")
	i2cAddr = flag.Int("i2c-addr", 0x76, "I2C connection address")
	i2cBus = flag.Int("i2c-bus", 1, "I2C connection bus line")
	//isDebug = flag.Bool("debug", false, "Debug mode for development")
)

// Environment variables
var (
	serviceURL = getEnv("RABBITMQ_SERVICE_URL", "rabbitmq-ha.default.svc.cluster.local")
	servicePort = getEnv("RABBITMQ_SERVICE_PORT", "5672")
	queue = getEnv("RABBITMQ_QUEUE", "sensor")
	exchange = getEnv("RABBITMQ_EXCHANGE", "sensor")
	secretPath = getEnv("SECRET_PATH", "/passwords")
)

var connection i2c.I2C

func main() {
	flag.Parse()
	setupSensor()
	defer connection.Close()

	// Start publishing messages
	ctx, done := context.WithCancel(context.Background())
	go func() {
		publish(redial(ctx, assembleURL()), read())
		done()
	}()
	<-ctx.Done()
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func assembleURL() string {
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
