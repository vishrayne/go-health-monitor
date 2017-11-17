package monit

import (
	"log"
	"path/filepath"

	toml "github.com/pelletier/go-toml"
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

//Stats representation
type Stats struct {
	CPU       *cpu       `json:"cpu"`
	Memory    *memory    `json:"memory"`
	Disk      *disk      `json:"disk"`
	Host      *host      `json:"host"`
	AccessLog *accessLog `json:"accessLog"`
	config    *toml.Tree
}

//Init does all the pre-requisites
func Init() *Stats {
	allStats := new(Stats)
	allStats.config = readConfigurations()
	return allStats
}

//CreateReport will create a new monit report
func (stats *Stats) CreateReport() {
	report := CreateSystemReport("monit")
	defer report.close()

	systemSummary := allStats(stats)
	report.writeSection("CPU", systemSummary.CPU.toPrettyJSON())
	report.writeSection("Memory", systemSummary.Memory.toPrettyJSON())
	report.writeSection("Disk", systemSummary.Disk.toPrettyJSON())
	report.writeSection("Host Information", systemSummary.Host.toPrettyJSON())
	report.writeSection("Access log", systemSummary.AccessLog.toPrettyJSON())
}

//AllStats collects everything to JSON
func (stats *Stats) AllStats() string {
	return asJSON(allStats(stats))
}

func allStats(stats *Stats) *Stats {
	stats.CPU = cpuStat(stats)
	stats.Memory = memoryStat(stats)
	stats.Disk = diskStat(stats)
	stats.Host = hostStat(stats)
	stats.AccessLog = accessLogSummary(stats)
	return stats
}

// CPUStat collects and returns cpu status as JSON
func (stats *Stats) CPUStat() string {
	return cpuStat(stats).toJSON()
}

func cpuStat(stats *Stats) *cpu {
	cpuDetails := newCPU()
	cpuDetails.collect()
	stats.CPU = cpuDetails
	return cpuDetails
}

// MemoryStat collects and returns memory status as JSON
func (stats *Stats) MemoryStat() string {
	return memoryStat(stats).toJSON()
}

func memoryStat(stats *Stats) *memory {
	memDetails := newMemory()
	memDetails.collect()
	stats.Memory = memDetails
	return memDetails
}

// DiskStat collects and returns disk status as JSON
func (stats *Stats) DiskStat() string {
	return diskStat(stats).toJSON()
}

func diskStat(stats *Stats) *disk {
	diskDetails := newDisk()
	diskDetails.collect()
	stats.Disk = diskDetails
	return diskDetails
}

// HostStat collects and returns host details as JSON
func (stats *Stats) HostStat() string {
	return hostStat(stats).toJSON()
}

func hostStat(stats *Stats) *host {
	hostDetails := newHost()
	hostDetails.collect()
	stats.Host = hostDetails
	return hostDetails
}

// AccessLogSummary summarizes the access log as JSON
func (stats *Stats) AccessLogSummary() string {
	return accessLogSummary(stats).toJSON()
}

func accessLogSummary(stats *Stats) *accessLog {
	accessLogParser := newAccessLogParser()
	accessConfig := stats.config.Get("access").(*toml.Tree)
	accessLogParser.parse(int(accessConfig.Get("line_count").(int64)), accessConfig.Get("file_name").(string), accessConfig.Get("show_log_entries").(bool))
	stats.AccessLog = accessLogParser
	return accessLogParser
}

func dealWithError(taskName string, err error) {
	if err != nil {
		log.Fatalf("%sTask failed: %v", taskName, err.Error())
	}
}

func readConfigurations() *toml.Tree {
	absPath, _ := filepath.Abs("config/settings.toml")
	config, _ := toml.LoadFile(absPath)
	return config
}
