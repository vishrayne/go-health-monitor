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
	defer report.close()

	report.writeSection("CPU", CPUStat())
	report.writeSection("Memory", MemoryStat())
	report.writeSection("Disk", DiskStat())
	report.writeSection("Host Information", HostStat())
	report.writeSection("Access log", AccessLogSummary())
}

// CPUStat collects and returns cpu status as JSON
func CPUStat() string {
	cpuDetails := newCPU()
	cpuDetails.collect()
	return cpuDetails.toJSON()
}

// MemoryStat collects and returns memory status as JSON
func MemoryStat() string {
	memDetails := newMemory()
	memDetails.collect()
	return memDetails.toJSON()
}

// DiskStat collects and returns disk status as JSON
func DiskStat() string {
	diskDetails := newDisk()
	diskDetails.collect()
	return diskDetails.toJSON()
}

// HostStat collects and returns host details as JSON
func HostStat() string {
	hostDetails := newHost()
	hostDetails.collect()
	return hostDetails.toJSON()
}

// AccessLogSummary summarizes the access log as JSON
func AccessLogSummary() string {
	accessLogParser := newAccessLogParser()
	accessLogParser.parse(100, "asset/Access-log-250917.txt", false)
	return accessLogParser.toJSON()
}

func dealWithError(taskName string, err error) {
	if err != nil {
		log.Fatalf("%sTask failed: %v", taskName, err.Error())
	}
}
