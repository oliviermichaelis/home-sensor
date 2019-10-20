package connect

import (
	"context"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"github.com/streadway/amqp"
	"log"
	"time"
)

// session composes an amqp.Connection with an amqp.Channel
type Session struct {
	*amqp.Connection
	*amqp.Channel
}


// Close tears the connection down, taking the channel with it.
func (s Session) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}

// redial continually connects to the URL, exiting the program when no longer possible/no route exists
func Redial(ctx context.Context, url string, exchange string, queue string) chan chan Session {
	sessions := make(chan chan Session)

	go func() {
		sess := make(chan Session)
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
				log.Printf("Attempted to connect to %s unsuccessfully. Retrying in: %v",environment.RabbitmqURL(),  duration)
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
			case sess <- Session{conn, ch}:
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
func Publish(sessions chan chan Session, messages <-chan []byte) {
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

// Consumes from a queue and sends the message to a channel used by the influxdb writer
func Subscribe(sessions chan chan Session, messages chan<- amqp.Delivery, tags <-chan uint64) {
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