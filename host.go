package monit

import (
	"runtime"
	"strings"

	sh "github.com/codeskyblue/go-sh"
)

type host struct {
	Uptime        string   `json:"uptime"`
	Kernal        string   `json:"kernal"`
	OS            string   `json:"os"`
	HostName      string   `json:"host_name"`
	Platform      string   `json:"platform"`
	LoggedInUsers []string `json:"logged_in_users"`
	LastReboot    string   `json:"last_reboot"`
}

func newHost() *host {
	return &host{}
}

func (h *host) toJSON() string {
	return asJSON(h)
}

func (h *host) collect() {
	h.fetchSystemInfo()
	h.fetchLoggedInUsers()
	h.fetchUptime()
}

func (h *host) fetchSystemInfo() {
	h.Kernal = parseString(sh.Command("uname", "-r").Output())
	h.OS = runtime.GOOS
	h.HostName = parseString(sh.Command("hostname").Output())
	h.Platform = parseString(sh.Command("lsb_release", "-d").Command("awk", "-F", "Description:", "{print $2}").Output())
	h.LastReboot = parseString(sh.Command("who", "-b").Command("awk", "{print $3,$4}").Output())
}

func (h *host) fetchUptime() {
	h.Uptime = parseString(sh.Command("uptime").Command("awk", "{print $3,$4}").Command("cut", "-f1", "-d,").Output())
}

func (h *host) fetchLoggedInUsers() {
	userString := parseString(sh.Command("who").Command("awk", "{print $1}").Command("sort").Command("uniq").Output())
	h.LoggedInUsers = strings.Fields(userString)
}
