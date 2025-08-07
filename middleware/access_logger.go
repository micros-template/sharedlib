package middleware

import (
	"encoding/json"
	"time"

	"10.1.20.130/dropping/log-management/pkg"
	ld "10.1.20.130/dropping/log-management/pkg/dto"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type userData struct {
	UserID string `json:"user_id"`
}

func AccessLogger(logEmitter pkg.LogEmitter, serviceName string, loger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log after response is sent
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()

		var userDataID string
		userDataHeader := c.Request.Header.Get("User-Data")
		if userDataHeader != "" {
			var ud userData
			err := json.Unmarshal([]byte(userDataHeader), &ud)
			if err == nil {
				userDataID = ud.UserID
			}
		}
		if userDataID == "" {
			userDataHeader = c.Writer.Header().Get("User-Data")
			if userDataHeader != "" {
				var ud userData
				err := json.Unmarshal([]byte(userDataHeader), &ud)
				if err == nil {
					userDataID = ud.UserID
				}
			}
		}

		logLevel := "INFO"
		if statusCode >= 400 && statusCode < 600 {
			logLevel = "ERR"
		}
		logData := map[string]interface{}{
			"type":       "access",
			"status":     statusCode,
			"method":     method,
			"path":       path,
			"ip":         clientIP,
			"user_agent": userAgent,
			"user_id":    userDataID,
			"latency":    duration.String(),
			"level":      logLevel,
		}
		logDataBytes, _ := json.Marshal(logData)
		logEmitter.EmitLog(c.Request.Context(), ld.LogMessage{
			Type:     logLevel,
			Service:  serviceName,
			Msg:      string(logDataBytes),
			Protocol: "HTTP",
		})
	}
}
