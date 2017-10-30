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

type stats struct {
	CPU              *cpu       `json:"cpu"`
	Memory           *memory    `json:"memory"`
	Disk             *disk      `json:"disk"`
	Host             *host      `json:"host"`
	AccessLogSummary *accessLog `json:"accessLog"`
}

//CreateReport will create a new monit report
func CreateReport() {
	report := CreateSystemReport("monit")
	defer report.close()

	systemSummary := allStats()
	report.writeSection("CPU", systemSummary.CPU.toPrettyJSON())
	report.writeSection("Memory", systemSummary.Memory.toPrettyJSON())
	report.writeSection("Disk", systemSummary.Disk.toPrettyJSON())
	report.writeSection("Host Information", systemSummary.Host.toPrettyJSON())
	report.writeSection("Access log", systemSummary.AccessLogSummary.toPrettyJSON())
}

// AllStats collects everything to JSON
func AllStats() string {
	return asJSON(allStats())
}

func allStats() *stats {
	allStats := new(stats)
	allStats.CPU = cpuStat()
	allStats.Memory = memoryStat()
	allStats.Disk = diskStat()
	allStats.Host = hostStat()
	allStats.AccessLogSummary = accessLogSummary()
	return allStats
}

// CPUStat collects and returns cpu status as JSON
func CPUStat() string {
	return cpuStat().toJSON()
}

func cpuStat() *cpu {
	cpuDetails := newCPU()
	cpuDetails.collect()
	return cpuDetails
}

// MemoryStat collects and returns memory status as JSON
func MemoryStat() string {
	return memoryStat().toJSON()
}

func memoryStat() *memory {
	memDetails := newMemory()
	memDetails.collect()
	return memDetails
}

// DiskStat collects and returns disk status as JSON
func DiskStat() string {
	return diskStat().toJSON()
}

func diskStat() *disk {
	diskDetails := newDisk()
	diskDetails.collect()
	return diskDetails
}

// HostStat collects and returns host details as JSON
func HostStat() string {
	return hostStat().toJSON()
}

func hostStat() *host {
	hostDetails := newHost()
	hostDetails.collect()
	return hostDetails
}

// AccessLogSummary summarizes the access log as JSON
func AccessLogSummary() string {
	return accessLogSummary().toJSON()
}

func accessLogSummary() *accessLog {
	accessLogParser := newAccessLogParser()
	accessLogParser.parse(100, "asset/Access-log-250917.txt", false)
	return accessLogParser
}

func dealWithError(taskName string, err error) {
	if err != nil {
		log.Fatalf("%sTask failed: %v", taskName, err.Error())
	}
}
