package auth

import (
	"backend/database"
	"backend/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetInfo 返回信息
// type responseJson struct {
// 	Id       int    `json:"id"`       // 用户Id
// 	Nickname string `json:"nickName"` // 用户昵称
// 	Type     int    `json:"password"` // 用户类型
// }

// @title             GetInfo
// @description       获取用户信息
// @auth              高宏宇         2022/2/12
// @param             c             请求句柄
func GetInfo(c *gin.Context) {
	//cookie, err := c.Cookie("camp-session")
	if _, err := c.Cookie("camp-session"); err != nil {
		whoAmIResponse := types.WhoAmIResponse{
			Code: types.LoginRequired,
			Data: types.TMember{},
		}
		c.JSON(http.StatusOK, whoAmIResponse)
		return
	}

	cookie, _ := c.Cookie("camp-session")

	dbsearchResult, user := database.GetUserInfoById(cookie)
	if dbsearchResult != "Success" {
		whoAmIResponse := types.WhoAmIResponse{
			Code: types.UserNotExisted,
			Data: types.TMember{},
		}
		c.JSON(http.StatusOK, whoAmIResponse)
		return
	}

	whoAmIResponse := types.WhoAmIResponse{
		Code: types.OK,
		Data: types.TMember{
			UserID:   strconv.Itoa(user.Id),
			Nickname: user.Nickname,
			Username: user.Name,
			UserType: user.Type,
		},
	}
	c.JSON(http.StatusOK, whoAmIResponse)
}
