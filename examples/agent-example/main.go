package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	machineID := registerMachine()
	log.Println("Machine registered with ID:", machineID)

	for {
		cpuPercent := collectCPU()
		ramPercent := collectRAM()
		pushMetrics(machineID, cpuPercent, ramPercent)
		time.Sleep(5 * time.Second)
	}
}

func registerMachine() int {
	body := map[string]string{
		"machine_name":  "my-laptop",
		"os":            "windows",
		"agent_version": "0.1.0",
	}
	b, _ := json.Marshal(body)
	resp, err := http.Post("http://localhost:8080/api/v1/machines/register", "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var res map[string]int
	json.NewDecoder(resp.Body).Decode(&res)
	return res["machine_id"]
}

func collectCPU() float64 {
	p, _ := cpu.Percent(0, false)
	return p[0]
}

func collectRAM() float64 {
	v, _ := mem.VirtualMemory()
	return v.UsedPercent
}

func pushMetrics(machineID int, cpuPercent, ramPercent float64) {
	body := map[string]interface{}{
		"machine_id":  machineID,
		"cpu_percent": cpuPercent,
		"ram_percent": ramPercent,
	}
	b, _ := json.Marshal(body)
	http.Post("http://localhost:8080/api/v1/metrics", "application/json", bytes.NewBuffer(b))
}
