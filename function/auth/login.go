package auth

import (
	"backend/database"
	"backend/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 请求信息
// type requestJson struct {
// 	Username string `form:"userName" loginRequest:"userName" binding:"required"` // 用户名
// 	Pssword  string `form:"password" loginRequest:"password" binding:"required"` // 密码
// }

// @title             Login
// @description       登录
// @auth              高宏宇         2022/2/12
// @param             c             请求句柄
func Login(c *gin.Context) {
	// 已经登录无需再次身份验证
	if cookie, err := c.Cookie("camp-session"); err == nil {
		loginResponse := types.LoginResponse{
			Code: http.StatusOK,
			Data: struct{ UserID string }{
				UserID: cookie,
			},
		}
		c.JSON(http.StatusOK, loginResponse)
		return
	}

	// 解析JSON
	var loginRequest types.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		loginResponse := types.LoginResponse{
			Code: http.StatusBadRequest,
			Data: struct{ UserID string }{
				UserID: "",
			},
		}
		c.JSON(http.StatusBadRequest, loginResponse)
		return
	}

	// 读取database获取user
	dbsearchResult, user := database.GetUserInfoByName(loginRequest.Username)
	if dbsearchResult != "Success" {
		loginResponse := types.LoginResponse{
			Code: types.UserNotExisted,
			Data: struct{ UserID string }{
				UserID: "",
			},
		}
		c.JSON(types.UserNotExisted, loginResponse)
		return
	}

	// 用户已删除
	if user.IsValid == 0 {
		loginResponse := types.LoginResponse{
			Code: types.UserHasDeleted,
			Data: struct{ UserID string }{
				UserID: "",
			},
		}
		c.JSON(types.UserHasDeleted, loginResponse)
		return
	}
	// 密码错误
	if loginRequest.Password != user.Password {
		loginResponse := types.LoginResponse{
			Code: types.WrongPassword,
			Data: struct{ UserID string }{
				UserID: "",
			},
		}
		c.JSON(types.WrongPassword, loginResponse)
		return
	}
	// 密码正确
	loginResponse := types.LoginResponse{
		Code: http.StatusOK,
		Data: struct{ UserID string }{
			UserID: user.Name,
		},
	}
	c.JSON(http.StatusOK, loginResponse)
	c.SetCookie("camp-session", user.Name, 3600, "/", "localhost", false, false)
	cookie, err1 := c.Cookie("camp-session")
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(cookie)
}
