package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @title             Logout
// @description       登出
// @auth              高宏宇         2022/2/12
// @param             c             请求句柄
func Logout(c *gin.Context) {
	var json requestJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 清除cookie
	c.SetCookie("camp-session", json.Username, -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
