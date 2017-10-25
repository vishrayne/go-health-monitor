package monit

import (
	"strings"

	sh "github.com/codeskyblue/go-sh"
)

type disk struct {
	Total          string `json:"total"`
	Used           string `json:"used"`
	Free           string `json:"free"`
	UsedPercent    string `json:"used_percent"`
	usedPercentInt float64
	Status         string `json:"status"`
}

func newDisk() *disk {
	return &disk{}
}

func (d *disk) toJSON() string {
	return asJSON(d)
}

func (d *disk) collect() {
	d.fetchUsage()
	d.checkStatus()
}

func (d *disk) fetchUsage() {
	// df -h --total | tail -1 | awk '{printf "total:%s used:%s free:%s used_percent:%s", $2, $3, $4, $5}'
	disk, err := sh.Command("df", "-h", "--total").Command("tail", "-1").Command("awk", "{printf \"total:%s used:%s free:%s used_percent:%s\", $2, $3, $4, $5}").Output()
	if err != nil {
		d.Total = "NA"
		d.Used = "NA"
		d.Free = "NA"
		d.UsedPercent = "NA"
		d.Status = "NA"
		return
	}

	stats := d.memStatsToMap(asString(disk))
	d.Total = stats["total"]
	d.Used = stats["used"]
	d.Free = stats["free"]
	d.UsedPercent = stats["used_percent"]
}

func (d *disk) memStatsToMap(stats string) map[string]string {
	if len(stats) <= 0 {
		return nil
	}

	parts := strings.Fields(stats)
	statMap := make(map[string]string, len(parts))
	for _, stat := range parts {
		pair := strings.Split(stat, ":")
		statMap[pair[0]] = asSafeString(pair[1])
	}

	return statMap
}

func (d *disk) checkStatus() {
	percent := strings.Replace(d.UsedPercent, "%", "", 1)
	d.usedPercentInt = asFloat(percent)

	switch {
	case d.usedPercentInt > 95:
		d.Status = Fatal
	case d.usedPercentInt > 75:
		d.Status = Caution
	case d.usedPercentInt > 50:
		d.Status = Warning
	case d.usedPercentInt <= 0:
		d.Status = "N.A"
	default:
		d.Status = Normal
	}
}
