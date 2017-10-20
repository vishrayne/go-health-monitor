package monit

import (
	"strings"

	sh "github.com/codeskyblue/go-sh"
)

type cpu struct {
	Count       int     `json:"cpus"`
	Utilization string  `json:"utilization"`
	Load1       float64 `json:"load1"`
	Load5       float64 `json:"load5"`
	Load15      float64 `json:"load15"`
}

func newCPU() *cpu {
	return &cpu{}
}

func (cpu *cpu) collectData() {
	cpu.Count = cpuCount()
	cpu.Utilization = cpuUtilization()

	avgLoad := averageLoad()
	cpu.Load1 = avgLoad["load1"]
	cpu.Load5 = avgLoad["load5"]
	cpu.Load15 = avgLoad["load15"]
}

func cpuCount() int {
	cpus, err := sh.Command("lscpu").Command("grep", "-e", "^CPU(s):").Command("cut", "-f2", "-d:").Command("awk", "{print $1}").Output()
	if err != nil {
		return -1
	}

	return asInteger(cpus)
}

func cpuUtilization() string {
	utilization, err := sh.Command("top", "-bn 2", "-d 0.01").Command("grep", "%Cpu").Command("tail", "-n 1").Command("awk", "{print $2+$4+$6}").Output()
	if err != nil {
		return "N.A"
	}

	return asString(utilization)
}

func averageLoad() map[string]float64 {
	load, err := sh.Command("uptime").Command("awk", "-F", "load average:", "{print $2}").Output()
	loadMap := make(map[string]float64, 3)
	var loadArray []string

	if err != nil {
		loadArray = []string{"-1", "-1", "-1"}
	} else {
		loadArray = strings.Split(string(load), ",")
		if len(loadArray) != 3 {
			loadArray = []string{"-1", "-1", "-1"}
		}
	}

	loadMap["load1"] = asFloat(loadArray[0])
	loadMap["load5"] = asFloat(loadArray[1])
	loadMap["load15"] = asFloat(loadArray[2])

	return loadMap
}
