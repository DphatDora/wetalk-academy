package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"wetalk-academy/package/logger"
	"wetalk-academy/web/admin"

	"github.com/gin-gonic/gin"
)

// ServeAdminLogsDashboard serves the static HTML UI for viewing logs.
func ServeAdminLogsDashboard(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", admin.LogsDashboardHTML)
}

// GetAdminLogFiles lists log files in the log directory (active + rotated backups).
func GetAdminLogFiles(c *gin.Context) {
	files, err := logger.ListLogFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"dir":   logger.LogDirectory(),
		"files": files,
	})
}

// GetAdminLogs returns the last N lines from a chosen log file as JSON.
// Query: file (base name, default: active log file), lines (default 200, max 2000).
func GetAdminLogs(c *gin.Context) {
	active := logger.LogFilePath()
	if active == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "log file path unavailable"})
		return
	}

	fileParam := c.Query("file")
	if fileParam == "" {
		fileParam = filepath.Base(active)
	}

	path, err := logger.ResolveLogFile(fileParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	n, _ := strconv.Atoi(c.DefaultQuery("lines", "200"))
	if n < 1 {
		n = 200
	}
	if n > 2000 {
		n = 2000
	}

	const maxRead = 2 * 1024 * 1024
	lines, err := logger.ReadTailLines(path, n, maxRead)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path":  path,
		"lines": lines,
	})
}
