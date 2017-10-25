package monit

import (
	"log"
)

const (
	// Normal status
	Normal = "Normal"
	// Warning status
	Warning = "Warning"
	// Caution status
	Caution = "Caution"
	// Fatal status
	Fatal = "Fatal"
)

//CreateReport will create a new monit report
func CreateReport() {
	report := CreateSystemReport("monit")

	cpuDetails := newCPU()
	cpuDetails.collect()
	report.writeSection("CPU", cpuDetails.toJSON())

	memDetails := newMemory()
	memDetails.collect()
	report.writeSection("Memory", memDetails.toJSON())

	diskDetails := newDisk()
	diskDetails.collect()
	report.writeSection("Disk", diskDetails.toJSON())

	hostDetails := newHost()
	hostDetails.collect()
	report.writeSection("Host Information", hostDetails.toJSON())

	report.close()
}

func dealWithError(taskName string, err error) {
	if err != nil {
		log.Fatalf("%sTask failed: %v", taskName, err.Error())
	}
}
