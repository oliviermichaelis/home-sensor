package main

import (
	"github.com/oliviermichaelis/home-sensor/pkg/connect"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"github.com/streadway/amqp"
	"log"
)

// Consumes from a queue and sends the message to a channel used by the influxdb writer
func subscribe(sessions chan chan connect.Session, messages chan<- amqp.Delivery, tags <-chan uint64) {
	queue := environment.GetQueue()
	for session := range sessions {
		sub := <-session

		if _, err := sub.QueueDeclare(queue, true, false, false, false, nil); err != nil {
			log.Printf("Cannot consume from exclusive queue: %q, %v", queue, err)
			return
		}

		deliveries, err := sub.Consume(queue, "", false, false, false, false, nil)
		if err != nil {
			log.Printf("Cannot consume from: %q, %v", queue, err)
			return
		}

		log.Printf("subscribed...")

		// Receives DeliveryTag as uint64 from writer and send ack to queue
		go func() {
			for tag := range tags {
				if err = sub.Ack(tag, false); err != nil {
					log.Fatal(err)
				}
			}
		}()

		// Sends message to writer
		for msg := range deliveries {
			messages <- msg
		}
	}
}
