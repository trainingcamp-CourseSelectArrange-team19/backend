package auth

import (
	"backend/types"

	"github.com/gin-gonic/gin"
)

// @title             AuthMiddleWare
// @description       基于cookie的身份验证中间件
// @auth              高宏宇         2022/2/12
// @return            func          gin.HandlerFunc类型函数
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
