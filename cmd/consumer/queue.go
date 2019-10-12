package main

import (
	"github.com/oliviermichaelis/home-sensor/pkg/connect"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"log"
)

// Consumes from a queue and sends the message to a channel used by the influxdb writer
func subscribe(sessions chan chan connect.Session, messages chan<- []byte) {
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

		for msg := range deliveries {
			messages <- []byte(msg.Body)
			if err = sub.Ack(msg.DeliveryTag, false); err != nil {		//TODO move ack to when insert to db succeeded
				log.Fatal(err)
			}
		}
	}
}
