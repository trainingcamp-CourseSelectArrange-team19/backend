package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	var json requestJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("camp-session", json.Username, -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
