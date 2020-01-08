package main

import (
	"log"
	"net/http"
)


func main() {
	//value := environment.SensorValues{
	//	Timestamp:   "20060102150405",
	//	Station:	"test",
	//	Temperature: 20.28,
	//	Humidity:    58.95,
	//	Pressure:    100615,
	//}
	//b, _ := json.Marshal(value)
	//ioutil.WriteFile("./binary_dump", b, 0644)
	http.HandleFunc("/measurements/climate", climateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
