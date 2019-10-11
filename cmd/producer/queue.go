package main

import (
	"github.com/oliviermichaelis/home-sensor/pkg/connect"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"github.com/streadway/amqp"
	"log"
)

// Publishes messages over a reconnecting session to a direct exchange.
// Messages are received over the channel 'messages'
func publish(sessions chan chan connect.Session, messages <-chan []byte) {
	for session := range sessions {
		var (
			running bool
			reading = messages
			pending = make(chan []byte, 1)
			confirm = make(chan amqp.Confirmation, 1)
		)

		pub := <-session

		// publisher confirms for this channel/connection
		if err := pub.Confirm(false); err != nil {
			log.Printf("publisher confirms not supported")
			close(confirm) // confirms not supported, simulate by always nacking
		} else {
			pub.NotifyPublish(confirm)
		}

		log.Printf("publishing...")

	Publish:
		for {
			var body []byte
			select {
			case confirmed, ok := <-confirm:
				if !ok {
					break Publish
				}
				if !confirmed.Ack {
					log.Printf("nack message %d, body: %q", confirmed.DeliveryTag, string(body))
				}
				reading = messages

			case body = <-pending:
				routingKey := "sensor"
				err := pub.Publish(environment.GetExchange(), routingKey, true, false, amqp.Publishing{
					Body: body,
				})
				// Retry failed delivery on the next session
				if err != nil {
					pending <- body
					if err = pub.Close(); err != nil {
						log.Fatal(err)
					}
					break Publish
				}

			case body, running = <-reading:
				// All messages consumed
				if !running {
					return
				}
				// Work on pending delivery until ack'd
				pending <- body
				reading = nil
			}
		}
	}
}
