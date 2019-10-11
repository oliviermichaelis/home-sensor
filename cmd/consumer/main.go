package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/oliviermichaelis/home-sensor/pkg/connect"
	"github.com/oliviermichaelis/home-sensor/pkg/environment"
	"io"
	"log"
	"os"
)

func main() {
	flag.Parse()

	// Start consuming messages
	ctx, done := context.WithCancel(context.Background())
	go func() {
		subscribe(connect.Redial(ctx, environment.AssembleURL(), environment.GetExchange(), environment.GetQueue()), write(os.Stdout))
		done()
	}()
	<-ctx.Done()
}

func write(w io.Writer) chan<- []byte {
	lines := make(chan []byte)
	go func() {
		for line := range lines {

			if _, err := fmt.Fprintln(w, string(line)); err != nil {
				log.Fatalln(err)
			}

		}
	}()
	return lines
}
