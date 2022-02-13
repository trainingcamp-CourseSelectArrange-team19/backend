package auth

import (
	"backend/database"
	"backend/tools"
	"backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type response struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickName"`
	Type     int    `json:"password"`
}

func GetInfo(c *gin.Context) {
	userName, _ := c.Cookie("camp-session")

	redisConn := redisPool.Get()
	defer redisConn.Close()
	var dbsearchResult string
	var user *database.User

	val, err := redisConn.Do("HMGET", hash, userName)
	if err == redis.Nil {
		dbsearchResult, user = database.GetUserInfoByName(userName)
		if dbsearchResult == "Success" {
			redisConn.Do("HSET", hash, userName, user)
		} else {
			c.JSON(types.UserNotExisted, gin.H{"status": types.UserNotExisted})
			return
		}
	} else if err != nil {
		tools.LogMsg(err)
		c.JSON(types.UnknownError, gin.H{"status": types.UnknownError})
		return
	} else {
		user = val.(*database.User)
	}

	json := response{
		Id:       user.Id,
		Nickname: user.Nickname,
		Type:     user.Type,
	}
	c.JSON(http.StatusOK, json)
}
