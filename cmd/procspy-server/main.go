package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/upayanmazumder/procspy/pkg/metrics"
	"github.com/upayanmazumder/procspy/pkg/store"
	"github.com/upayanmazumder/procspy/pkg/websocket"
)

var (
	s   *store.Store
	hub = websocket.NewHub()
)

func main() {
	s = store.NewStore("postgres://procspy:procspy@localhost:5432/procspy?sslmode=disable")
	metrics.RegisterMetrics()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/machines/register", registerMachineHandler).Methods("POST")
	r.HandleFunc("/api/v1/metrics", pushMetricsHandler).Methods("POST")
	r.HandleFunc("/ws", websocket.ServeWS(hub))
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("procspy server live"))
	})

	log.Println("procspy server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

type RegisterRequest struct {
	MachineName  string `json:"machine_name"`
	OS           string `json:"os"`
	AgentVersion string `json:"agent_version"`
}

type MetricsRequest struct {
	MachineID int     `json:"machine_id"`
	CPU       float64 `json:"cpu_percent"`
	RAM       float64 `json:"ram_percent"`
}

// Handler to register a new machine
func registerMachineHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := s.RegisterMachine(req.MachineName, req.OS, req.AgentVersion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"machine_id": id})
}

// Handler to push metrics from agent
func pushMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var req MetricsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save to DB
	if err := s.SaveMetric(req.MachineID, req.CPU, req.RAM); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update Prometheus metrics
	metrics.CPUPercent.WithLabelValues(strconv.Itoa(req.MachineID)).Set(req.CPU)
	metrics.RAMPercent.WithLabelValues(strconv.Itoa(req.MachineID)).Set(req.RAM)

	// Broadcast to WebSocket clients
	msg, _ := json.Marshal(req)
	hub.Broadcast(msg)

	json.NewEncoder(w).Encode(map[string]string{"status": "metrics_received"})
}
