package main

import (
	"flag"
	"fmt"
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
	flagLocal := flag.Bool("local", false, "Enable local setup")
	flagDebug := flag.Bool("debug", false, "Enable debug output")
	flagInfluxURL := flag.String("influx.url", "localhost", "URL of influxdb database")
	flagInfluxPort := flag.Int("influx.port", 9999, "Influxdb port")
	//flagInfluxUsernamePath := flag.String("influx.secrets.username", "/secrets/influx/username", "Path at which username is stored")
	//flagInfluxPasswordPath := flag.String("influx.secrets.password", "/secrets/influx/password", "Path at which password is stored")
	flagInfluxToken := flag.String("influx.secrets.token", "/secrets/influx/token", "Path at which token is stored")
	flag.Parse()

	var (
		err error
		//influxUser = "test"
		//influxPassword = "test"
		influxToken = "test"
	)

	if !*flagLocal {
		//influxUser, err = readSecret(*flagInfluxUsernamePath)
		//if err != nil {
		//	log.Fatalf("can't read username: %v", err)
		//}
		//
		//influxPassword, err = readSecret(*flagInfluxPasswordPath)
		//if err != nil {
		//	log.Fatalf("can't read password: %v", err)
		//}
		if influxToken, err = readSecret(*flagInfluxToken); err != nil {
			log.Fatalf("can't read token: %v", err)
		}

	}

	// TODO add prometheus /metrics
	// TODO add profiler /pprof
	mux := http.NewServeMux()
	repo := apiserver.NewInfluxDbRepository(*flagDebug, fmt.Sprintf("http://%s:%d", *flagInfluxURL, *flagInfluxPort), influxToken)
	server := apiserver.NewServer(*flagDebug, mux, ":8080", repo)
	log.Fatal(server.Start())
}
