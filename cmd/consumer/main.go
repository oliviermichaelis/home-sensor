package main

import (
	"context"
	"encoding/json"
	"flag"
	influxdb "github.com/influxdata/influxdb1-client/v2"
	"github.com/oliviermichaelis/home-sensor/pkg/connect"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"github.com/oliviermichaelis/home-sensor/pkg/healthcheck"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	flag.Parse()

	// Connect to influxdb
	client := setupClient()
	healthcheck.Health.SetStatus(healthcheck.ServiceInfluxDB, true)
	defer client.Close()
	createDatabase(client)

	// Start HTTP Server for healthchecks
	healthcheck.Server(healthcheck.Consumer)

	// Start consuming messages
	ctx, done := context.WithCancel(context.Background())
	go func() {
		tags := make(chan uint64)
		connect.Subscribe(connect.Redial(ctx, environment.RabbitmqURL(), environment.GetExchange(), environment.GetQueue()), write(client, tags), tags)
		done()
	}()
	<-ctx.Done()
}

func write(client influxdb.Client, tags chan<- uint64) chan<- amqp.Delivery {
	consume := make(chan amqp.Delivery)
	go func() {
		for consumed := range consume {
			values := environment.SensorValues{}
			if err := json.Unmarshal(consumed.Body, &values); err != nil {
				log.Fatalf("Could not unmarshal json! %v", err)
			}
			if err := insertPoint(client, values); err != nil {
				log.Fatal(err)
			}
			tags <- consumed.DeliveryTag
		}
	}()
	return consume
}
