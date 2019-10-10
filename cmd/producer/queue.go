package main

import (
	"context"
	"github.com/streadway/amqp"
	"log"
	"time"
)

// session composes an amqp.Connection with an amqp.Channel
type session struct {
	*amqp.Connection
	*amqp.Channel
}

// Close tears the connection down, taking the channel with it.
func (s session) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}

// redial continually connects to the URL, exiting the program when no longer possible/no route exists
func redial(ctx context.Context, url string) chan chan session {
	sessions := make(chan chan session)

	go func() {
		sess := make(chan session)
		defer close(sessions)

		for {
			select {
			case sessions <- sess:
			case <-ctx.Done():
				log.Println("shutting down session factory")
				return
			}

			var conn *amqp.Connection
			var err error
			for i := 0; i < 3; i++ {
				conn, err = amqp.Dial(url)
				if err == nil {
					break
				}
				duration := time.Duration(i) * 10 * time.Second
				log.Println("Attempted to connect to RabbitMQ unsuccessfully. Retrying in: ", duration)
				time.Sleep(duration)
			}
			if err != nil {
				log.Fatalf("Cannot create channel: %v", err)
			}
			log.Printf("Connected successfully to %s", url)

			/*
			conn, err := amqp.Dial(url)
			if err != nil {
				log.Fatalf("cannot (re)dial: %v: %q", err, url)
			}
			*/

			ch, err := conn.Channel()
			if err != nil {
				log.Fatalf("cannot create channel: %v", err)
			}

			if err := ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); err != nil {
				log.Fatalf("cannot declare direct exchange: %v", err)
			}

			_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
			if err != nil {
				log.Fatal(err)
			}

			if err := ch.QueueBind(queue, "sensor", exchange, false, nil ); err != nil {
				log.Fatalf("cannot bind queue to exchange: %v", err)
			}

			select {
			case sess <- session{conn, ch}:
			case <-ctx.Done():
				log.Println("shutting down new session")
				return
			}
		}
	}()

	return sessions
}

// Publishes messages over a reconnecting session to a direct exchange.
// Messages are received over the channel 'messages'
func publish(sessions chan chan session, messages <-chan []byte) {
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
				err := pub.Publish(exchange, routingKey, true, false, amqp.Publishing{
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
