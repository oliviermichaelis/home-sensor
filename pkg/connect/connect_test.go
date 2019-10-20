package connect

import (
	"context"
	"encoding/json"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"github.com/streadway/amqp"
	"log"
	"reflect"
	"testing"
)

func reader(value environment.SensorValues) <-chan []byte {
	lines := make(chan []byte)
	go func() {
		defer close(lines)
		j, err := json.Marshal(value)
		if err != nil {
			log.Fatal(err)
		}
		lines <- j
	}()
	return lines
}

func writer(received chan<- environment.SensorValues, tags chan<- uint64) chan<- amqp.Delivery {
	consume := make(chan amqp.Delivery)
	go func() {
		consumed := <- consume
		var value = environment.SensorValues{}
		if err := json.Unmarshal(consumed.Body, &value); err != nil {
			log.Fatalf("Could not unmarshal json! %v", err)
		}
		tags <- consumed.DeliveryTag
		received <- value
	}()

	return consume
}

func TestConnect(t *testing.T) {
	value := environment.SensorValues{
		Timestamp:   "20060102150405",
		Temperature: 20.28,
		Humidity:    58.95,
		Pressure:    100615,
	}

	ctx, done := context.WithCancel(context.Background())
	Publish(Redial(ctx, environment.RabbitmqURL(), environment.GetExchange(), environment.GetQueue()), reader(value))
	done()

	tags := make(chan uint64)
	receive := make(chan environment.SensorValues)
	ctx, done = context.WithCancel(context.Background())
	go func() {
		Subscribe(Redial(ctx, environment.RabbitmqURL(), environment.GetExchange(), environment.GetQueue()), writer(receive, tags), tags)
		done()
	}()
	consumed := <- receive

	if !reflect.DeepEqual(value, consumed) {
		t.Error("Published value doesn't correspond to received value!")
	}
}
