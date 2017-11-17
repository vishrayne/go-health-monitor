package main

import monit "github.com/vishrayne/go-monit"

func main() {
	stats := monit.Init()
	stats.CreateReport()
}
