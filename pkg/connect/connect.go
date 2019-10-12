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