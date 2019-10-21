package healthcheck

import (
	"encoding/json"
	"log"
	"net/http"
)

type Health struct {
	Status string
}

func (h Health) SetPass() {
	status.Status = "pass"
}

func (h Health) SetWarn() {
	status.Status = "warn"
}

func (h Health) SetFail() {
	status.Status = "fail"
}

var status = Health{Status: "fail"}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(status)
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

func Server() {
	log.Println("Starting healthcheck endpoint")
	http.HandleFunc("/health", healthHandler)
	go func() {
		log.Fatal(http.ListenAndServe(":80", nil))
	}()
}
