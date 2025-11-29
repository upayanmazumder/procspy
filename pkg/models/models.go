package models

import "time"

type Machine struct {
	ID           int       `json:"id"`
	MachineName  string    `json:"machine_name"`
	OS           string    `json:"os"`
	AgentVersion string    `json:"agent_version"`
	RegisteredAt time.Time `json:"registered_at"`
}

type Metric struct {
	ID        int       `json:"id"`
	MachineID int       `json:"machine_id"`
	CPU       float64   `json:"cpu_percent"`
	RAM       float64   `json:"ram_percent"`
	Collected time.Time `json:"collected_at"`
}
