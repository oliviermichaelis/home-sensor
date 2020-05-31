package main

import (
	"flag"
	"github.com/oliviermichaelis/home-sensor/pkg/apiserver"
	"io/ioutil"
	"log"
	"net/http"
)

func readSecret(path string) (string, error) {
	s, err := ioutil.ReadFile(path)
	return string(s), err
}

func main() {
	flagDebug := flag.Bool("debug", false, "Enable debug output")
	flagInfluxURL := flag.String("influx.url", "localhost", "URL of influxdb database")
	flagInfluxPort := flag.Int("influx.port", 8086, "Influxdb port")
	flagInfluxUsernamePath := flag.String("influx.secrets.username", "/secrets/influx/username", "Path at which username is stored")
	flagInfluxPasswordPath := flag.String("influx.secrets.password", "/secrets/influx/password", "Path at which password is stored")
	flag.Parse()

	influxUser, err := readSecret(*flagInfluxUsernamePath)
	if err != nil {
		log.Fatalf("can't read username: %v", err)
	}

	influxPassword, err := readSecret(*flagInfluxPasswordPath)
	if err != nil {
		log.Fatalf("can't read password: %v", err)
	}

	// TODO add prometheus /metrics
	// TODO add profiler /pprof
	mux := http.NewServeMux()
	repo := apiserver.NewInfluxDbRepository(*flagDebug, *flagInfluxURL, *flagInfluxPort, influxUser, influxPassword)
	server := apiserver.NewServer(*flagDebug, mux, ":8080", repo)
	log.Fatal(server.Start())
}
