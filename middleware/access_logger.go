package middleware

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type userData struct {
	UserID string `json:"user_id"`
}

func AccessLogger(loger zerolog.Logger) gin.HandlerFunc {
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

		loger.Info().
			Str("type", "access").
			Int("status", statusCode).
			Str("method", method).
			Str("path", path).
			Str("ip", clientIP).
			Str("user_agent", userAgent).
			Str("user_id", userDataID).
			Dur("latency", duration).
			Msg("incoming request")
	}
}
