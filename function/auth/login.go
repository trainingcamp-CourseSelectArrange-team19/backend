package auth

import (
	"backend/database"
	"backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type request struct {
	Username string `form:"userName" json:"userName" binding:"required"`
	Pssword  string `form:"password" json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	if _, err := c.Cookie("camp-session"); err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "200"})
		return
	}

	var json request
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err, user := database.GetUserInfoByName(json.Username)

	status := http.StatusOK
	if err != "Success" {
		status = types.UserNotExisted
	} else if !user.IsValid {
		status = types.UserHasDeleted
	} else if json.Pssword != user.Password {
		status = types.WrongPassword
	}
	c.JSON(status, gin.H{"status": status})

	if status == http.StatusOK {
		c.SetCookie("camp-session", json.Username, 3600, "/", "localhost", false, true)
	}
}
