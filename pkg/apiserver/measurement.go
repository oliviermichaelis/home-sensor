package apiserver

import (
	"time"
)

type measurement struct {
	Timestamp   time.Time `json:"timestamp"`
	Station     string    `json:"station"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Pressure    float64   `json:"pressure"`
}
