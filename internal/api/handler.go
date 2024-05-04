package api

import (
	"benzinga-webhook-receiver/pkg/config"
	"benzinga-webhook-receiver/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LogEntry struct {
	UserID    int                    `json:"user_id"`
	Total     float64                `json:"total"`
	Title     string                 `json:"title"`
	Meta      map[string]interface{} `json:"meta"`
	Completed bool                   `json:"completed"`
}

type LoginRecord struct {
	Time string `json:"time"`
	IP   string `json:"ip"`
}

// Helper function to parse and validate time format
func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

func SetupRoutes(router *gin.Engine) {
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.POST("/log", HandleLogPost)
}

func HandleLogPost(c *gin.Context) {
	var entry models.LogEntry
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON data"})
		return
	}

	// Check required fields
	if entry.UserID <= 0 || entry.Total <= 0 || entry.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid 'user_id', 'total', or 'title'"})
		return
	}

	// Validate logins: each login must have a valid time and IP
	for _, login := range entry.Meta.Logins {
		formattedTime := login.Time.Format(time.RFC3339)
		_, err := time.Parse(time.RFC3339, formattedTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format"})
			return
		}
		if login.IP == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing IP address"})
			return
		}
	}

	// Validate phone numbers
	if entry.Meta.PhoneNumbers.Home == "" || entry.Meta.PhoneNumbers.Mobile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing phone numbers"})
		return
	}

	// Add log entry to batch
	config.AddToBatch(entry)
	c.JSON(http.StatusOK, gin.H{"status": "logged"})
}
