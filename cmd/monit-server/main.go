package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	monit "github.com/vishrayne/go-monit"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.GET("/", rootHandler)
	engine.GET("/ping", pingHandler)
	engine.GET("/summary", reportHandler)
	engine.GET("/cpu", cpuHandler)
	engine.GET("/memory", memoryHandler)
	engine.GET("/disk", diskHandler)
	engine.GET("/host", hostHandler)
	engine.GET("/accesslog", accessLogHandler)

	engine.Run(":8080")
}

func rootHandler(c *gin.Context) {
	pingHandler(c)
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Service is live!", "title": "lvb-system-monit"})
}

func reportHandler(c *gin.Context) {
	c.JSON(http.StatusOK, monit.AllStats())
}

func cpuHandler(c *gin.Context) {
	c.JSON(http.StatusOK, monit.CPUStat())
}

func memoryHandler(c *gin.Context) {
	c.JSON(http.StatusOK, monit.MemoryStat())
}

func diskHandler(c *gin.Context) {
	c.JSON(http.StatusOK, monit.DiskStat())
}

func hostHandler(c *gin.Context) {
	c.JSON(http.StatusOK, monit.HostStat())
}

func accessLogHandler(c *gin.Context) {
	c.JSON(http.StatusOK, monit.AccessLogSummary())
}
