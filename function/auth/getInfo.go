package auth

import (
	"backend/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickName"`
	Type     string `json:"password"`
}

func GetInfo(c *gin.Context) {
	userName, _ := c.Cookie("camp-session")

	user := database.GetUserInfoByName(userName)

	json := response{
		Id:       user.Id,
		Nickname: user.Nickname,
		Type:     user.Type,
	}
	c.JSON(http.StatusOK, json)
}
