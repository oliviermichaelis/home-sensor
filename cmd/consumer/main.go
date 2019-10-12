package main

import (
	"context"
	"encoding/json"
	"flag"
	influxdb "github.com/influxdata/influxdb1-client/v2"
	"github.com/oliviermichaelis/home-sensor/pkg/connect"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"log"
)

func main() {
	flag.Parse()

	// Connect to influxdb
	client := setupClient()
	defer client.Close()
	createDatabase(client)

	// Start consuming messages
	ctx, done := context.WithCancel(context.Background())
	go func() {
		subscribe(connect.Redial(ctx, environment.RabbitmqURL(), environment.GetExchange(), environment.GetQueue()), write(client))
		done()
	}()
	<-ctx.Done()
}

func write(client influxdb.Client) chan<- []byte {
	consume := make(chan []byte)
	go func() {
		for consumed := range consume {
			values := environment.SensorValues{}
			if err := json.Unmarshal(consumed, &values); err != nil {
				log.Fatalf("Could not unmarshal json! %v", err)
			}
			insertPoint(client, values)
		}
	}()
	return consume
}
