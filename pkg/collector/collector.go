package collector

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func CollectCPU() (float64, error) {
	percents, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	return percents[0], nil
}

func CollectRAM() (float64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return v.UsedPercent, nil
}
