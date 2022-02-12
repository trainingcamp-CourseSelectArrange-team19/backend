package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type request struct {
	Username string `form:"userName" json:"userName" binding:"required"`
	Pssword  string `form:"password" json:"password" binding:"required"`
}

func Logout(c *gin.Context) {
	var json request
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("camp-session", json.Username, -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"status": "200"})
}
