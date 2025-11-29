package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	CPUPercent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "procspy_cpu_percent",
			Help: "CPU usage percentage per machine",
		},
		[]string{"machine_id"},
	)

	RAMPercent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "procspy_ram_percent",
			Help: "RAM usage percentage per machine",
		},
		[]string{"machine_id"},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(CPUPercent)
	prometheus.MustRegister(RAMPercent)
}
