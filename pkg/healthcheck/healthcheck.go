package healthcheck

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

const ServiceRabbitMQ = "rabbitmq"
const ServiceInfluxDB = "influxdb"
const statusPass = "pass"
const statusFail = "fail"
const statusWarn = "warn"
const Producer = 0
const Consumer = 1

var Health = health{Status: "fail", Details: []service{},}
var isConsumer = true

type health struct {
	sync.RWMutex
	Status string
	Details []service
}

type service struct {
	Name string
	Status bool
}

func (h *health) setPass() {
	h.Status = statusPass
}

func (h *health) setWarn() {
	h.Status = statusWarn
}

func (h *health) setFail() {
	h.Status = statusFail
}

func (h *health) checkStatus() {
	rStatus, iStatus := false, false
	for _, service := range h.Details {
		if service.Name == ServiceRabbitMQ {
			rStatus = service.Status
		} else if service.Name == ServiceInfluxDB {
			iStatus = service.Status
		}
	}

	// The variable isConsumer is set when the Server is started. If it's the Producer, only rStatus needs to be true
	// for the expression to be evaluated to true
	if rStatus && (!isConsumer || iStatus) {
		h.setPass()
	} else {
		h.setFail()
	}
}

func (h *health) SetStatus(serviceName string, status bool) {
	h.Lock()
	defer h.Unlock()
	for _, service := range h.Details {
		if service.Name == serviceName {
			service.Status = status
			h.checkStatus()
			break
		}
	}

	h.Details = append(h.Details, service{
		Name:   serviceName,
		Status: status,
	})
	h.checkStatus()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	Health.Lock()
	j, err := json.Marshal(Health)
	if Health.Status == statusFail {
		w.WriteHeader(http.StatusNotFound)
	}
	Health.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(j); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

}

func Server(service int) {
	if service == Producer {
		isConsumer = false
	}

	log.Println("Starting healthcheck endpoint")
	http.HandleFunc("/health", healthHandler)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}
