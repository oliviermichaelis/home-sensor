package main

import (
	"context"
	"flag"
	"github.com/oliviermichaelis/home-sensor/pkg/connect"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"github.com/oliviermichaelis/home-sensor/pkg/healthcheck"
	"time"
)

// cli flags
var (
	samplerate = flag.Duration("samplerate", 60 * time.Second, "Sample rate at which sensor is polled")
	i2cAddr = flag.Int("i2c-addr", 0x76, "I2C connection address")
	i2cBus = flag.Int("i2c-bus", 1, "I2C connection bus line")
)

func main() {
	flag.Parse()
	connection := setupSensor()
	defer connection.Close()

	// Start HTTP Server for healthchecks
	healthcheck.Server(healthcheck.Producer)

	// Start publishing messages
	ctx, done := context.WithCancel(context.Background())
	go func() {
		connect.Publish(connect.Redial(ctx, environment.RabbitmqURL(), environment.GetExchange(), environment.GetQueue()), read())
		done()
	}()
	<-ctx.Done()
}
