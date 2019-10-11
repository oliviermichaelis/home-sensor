package main

import (
	"context"
	"github.com/oliviermichaelis/home-sensor/pkg/connect"
	"os"
	"testing"
)

var ctx context.Context
var done context.CancelFunc

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	ctx, done = context.WithCancel(context.Background())
}

func shutdown() {
	done()
	<-ctx.Done()
}

func mockedReader() <-chan []byte {
	lines := make(chan []byte)
	go func() {
		defer close(lines)
		lines <- []byte("test")
	}()
	return lines
}

func TestPublishSuccessful(t *testing.T) {
	url := "amqp://guest:guest@127.0.0.1:5672/"
	publish(connect.Redial(ctx, url, "sensor", "sensor"), mockedReader())
}
