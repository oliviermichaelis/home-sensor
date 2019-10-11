package main

import (
	"context"
	"flag"
	"github.com/d2r2/go-i2c"
	"github.com/oliviermichaelis/home-sensor/pkg/connect"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"time"
)

// cli flags
var (
	samplerate = flag.Duration("samplerate", 60 * time.Second, "Sample rate at which sensor is polled")
	i2cAddr = flag.Int("i2c-addr", 0x76, "I2C connection address")
	i2cBus = flag.Int("i2c-bus", 1, "I2C connection bus line")
)

var connection i2c.I2C

func main() {
	flag.Parse()
	setupSensor()
	defer connection.Close()

	// Start publishing messages
	ctx, done := context.WithCancel(context.Background())
	go func() {
		publish(connect.Redial(ctx, environment.AssembleURL(), environment.GetExchange(), environment.GetQueue()), read())
		done()
	}()
	<-ctx.Done()
}
