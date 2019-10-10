package main

import (
	"context"
	"flag"
	"github.com/d2r2/go-i2c"
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
	url = getEnv("AMQP_URL", "amqp://guest:guest@127.0.0.1:5672/")
	queue = getEnv("RABBITMQ_QUEUE", "sensor")
	exchange = getEnv("RABBITMQ_EXCHANGE", "sensor")
)

var connection i2c.I2C

func main() {
	flag.Parse()
	setupSensor()
	defer connection.Close()

	// Start publishing messages
	ctx, done := context.WithCancel(context.Background())
	go func() {
		publish(redial(ctx, url), read())
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
