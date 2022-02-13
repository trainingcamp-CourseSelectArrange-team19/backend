package auth

import (
	"backend/types"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := c.Cookie("camp-session"); err == nil {
			c.Next()
			return
		}
		c.JSON(types.LoginRequired, gin.H{"status": types.LoginRequired})
		c.Abort()
	}
}
