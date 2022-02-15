package auth

import (
	"backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @title             Logout
// @description       登出
// @auth              高宏宇         2022/2/12
// @param             c             请求句柄
func Logout(c *gin.Context) {
	if _, err := c.Cookie("camp-session"); err != nil {
		logoutResponse := types.LogoutResponse{
			Code: types.LoginRequired,
		}
		c.JSON(types.LoginRequired, logoutResponse)
		return
	}

	var logoutRequest types.LogoutRequest
	if err := c.ShouldBindJSON(&logoutRequest); err != nil {
		logoutResponse := types.LogoutResponse{
			Code: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, logoutResponse)
		return
	}

	// 清除cookie
	c.SetCookie("camp-session", "", -1, "/", "localhost", false, false)

	logoutResponse := types.LogoutResponse{
		Code: http.StatusOK,
	}
	c.JSON(http.StatusOK, logoutResponse)
}
