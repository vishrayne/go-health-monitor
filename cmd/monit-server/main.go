package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/gin"
	monit "github.com/vishrayne/go-monit"
)

const monitStatus string = "monit_stats"

func main() {

	// gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	engine.Use(monitMiddleware())
	engine.Use(secure.Secure(secure.Options{
		AllowedHosts:          []string{"localhost:8080", "google.com"},
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	}))
	engine.Use(cors.Default())

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

func monitMiddleware() gin.HandlerFunc {
	// one-time initialization
	stats := monit.Init()

	return func(c *gin.Context) {
		c.Set(monitStatus, stats)
		c.Next()
	}
}

func rootHandler(c *gin.Context) {
	pingHandler(c)
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Service is live!", "title": "lvb-system-monit"})
}

func reportHandler(c *gin.Context) {
	stats := c.MustGet(monitStatus).(*monit.Stats)
	c.JSON(http.StatusOK, stats.AllStats())
}

func cpuHandler(c *gin.Context) {
	stats := c.MustGet(monitStatus).(*monit.Stats)
	c.JSON(http.StatusOK, stats.CPUStat())
}

func memoryHandler(c *gin.Context) {
	stats := c.MustGet(monitStatus).(*monit.Stats)
	c.JSON(http.StatusOK, stats.MemoryStat())
}

func diskHandler(c *gin.Context) {
	stats := c.MustGet(monitStatus).(*monit.Stats)
	c.JSON(http.StatusOK, stats.DiskStat())
}

func hostHandler(c *gin.Context) {
	stats := c.MustGet(monitStatus).(*monit.Stats)
	c.JSON(http.StatusOK, stats.HostStat())
}

func accessLogHandler(c *gin.Context) {
	stats := c.MustGet(monitStatus).(*monit.Stats)
	c.JSON(http.StatusOK, stats.AccessLogSummary())
}
