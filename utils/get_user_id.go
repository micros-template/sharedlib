package utils

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type UserData struct {
	UserId string `json:"user_id"`
}

func GetUserId(c *gin.Context) string {
	userDataHeader := c.Request.Header.Get("User-Data")
	if userDataHeader != "" {
		var ud UserData
		err := json.Unmarshal([]byte(userDataHeader), &ud)
		if err == nil {
			return ud.UserId
		}
	}
	return ""
}
