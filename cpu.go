package monit

import (
	"log"
	"strings"

	sh "github.com/codeskyblue/go-sh"
)

const (
	//SafeLoadPercent is the maximum allowed value for average load
	SafeLoadPercent = 0.7
)

type cpu struct {
	Count           int     `json:"cpus"`
	Utilization     string  `json:"utilization"`
	Load1           float64 `json:"load1"`
	Load5           float64 `json:"load5"`
	Load15          float64 `json:"load15"`
	MaxAllowedLoad  float64 `json:"max_allowed_load"`
	Status          string  `json:"status"`
	isLoadAvailable bool
}

func newCPU() *cpu {
	return &cpu{}
}

func (cpu *cpu) toJSON() string {
	return asJSON(cpu)
}

func (cpu *cpu) collectData() {
	cpu.fetchCPUCount()
	cpu.fetchCPUUtilization()
	cpu.fetchAverageLoad()
	cpu.checkStatus()
}

func (cpu *cpu) fetchCPUCount() {
	cpus, err := sh.Command("lscpu").Command("grep", "-e", "^CPU(s):").Command("cut", "-f2", "-d:").Command("awk", "{print $1}").Output()
	if err != nil {
		cpu.Count = -1
	} else {
		cpu.Count = asInteger(cpus)
	}
}

func (cpu *cpu) fetchCPUUtilization() {
	utilization, err := sh.Command("top", "-bn 2", "-d 0.01").Command("grep", "%Cpu").Command("tail", "-n 1").Command("awk", "{print $2+$4+$6}").Output()
	if err != nil {
		cpu.Utilization = "N.A"
	} else {
		cpu.Utilization = asString(utilization)
	}
}

func (cpu *cpu) fetchAverageLoad() {
	load, err := sh.Command("uptime").Command("awk", "-F", "load average:", "{print $2}").Output()
	var loadArray []string

	if err != nil {
		loadArray = []string{"-1", "-1", "-1"}
	} else {
		loadArray = strings.Split(string(load), ",")
		if len(loadArray) != 3 {
			loadArray = []string{"-1", "-1", "-1"}
		}
	}

	cpu.Load1 = asFloat(loadArray[0])
	cpu.Load5 = asFloat(loadArray[1])
	cpu.Load15 = asFloat(loadArray[2])
	cpu.isLoadAvailable = true
}

func (cpu *cpu) checkStatus() {
	if !cpu.isLoadAvailable {
		log.Println("Avg. load is not available. Fetching...")
		cpu.fetchAverageLoad()
	}

	healthyLimit := SafeLoadPercent * float64(cpu.Count)
	cpu.MaxAllowedLoad = healthyLimit

	switch {
	case cpu.Load15 >= healthyLimit:
		cpu.Status = Fatal
	case cpu.Load5 >= healthyLimit:
		cpu.Status = Caution
	case cpu.Load1 >= healthyLimit:
		cpu.Status = Warning
	default:
		cpu.Status = Normal
	}

}
