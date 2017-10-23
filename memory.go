package monit

import (
	"strconv"
	"strings"

	sh "github.com/codeskyblue/go-sh"
)

type memory struct {
	Total          string `json:"total"`
	totalInt       uint64
	Used           string `json:"used"`
	usedInt        uint64
	Free           string `json:"free"`
	freeInt        uint64
	UsedPercent    string `json:"used_percent"`
	usedPercentInt float64
	Status         string `json:"status"`
}

func newMemory() *memory {
	return &memory{}
}

func (m *memory) toJSON() string {
	return asJSON(m)
}

func (m *memory) collect() {
	m.fetchUsage()
}

func (m *memory) fetchUsage() {
	mem, err := sh.Command("free", "-b").Command("head", "-2").Command("tail", "-1").Command("awk", "{printf \"total:%s used:%s free:%s\",  $2, $3, $4}").Output()
	if err != nil {
		m.Total = "NA"
		m.Used = "NA"
		m.Free = "NA"
		m.UsedPercent = "NA"
		m.Status = "NA"
		return
	}

	stats := m.memStatsToMap(asString(mem))
	m.totalInt = stats["total"]
	m.usedInt = stats["used"]
	m.freeInt = stats["free"]

	m.Total = asHumanBytes(m.totalInt)
	m.Used = asHumanBytes(m.usedInt)
	m.Free = asHumanBytes(m.freeInt)

	m.usedPercentInt = (float64(m.usedInt) / float64(m.totalInt)) * 100
	m.UsedPercent = strconv.FormatFloat(m.usedPercentInt, 'f', 2, 64)
}

func (m *memory) memStatsToMap(stats string) map[string]uint64 {
	if len(stats) <= 0 {
		return nil
	}

	parts := strings.Fields(stats)
	statMap := make(map[string]uint64, len(parts))
	for _, stat := range parts {
		pair := strings.Split(stat, ":")
		statMap[pair[0]] = asUInt64(pair[1])
	}

	return statMap
}
