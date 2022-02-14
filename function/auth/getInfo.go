package auth

import (
	"backend/database"
	"backend/tools"
	"backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// GetInfo 返回信息
type responseJson struct {
	Id       int    `json:"id"`       // 用户Id
	Nickname string `json:"nickName"` // 用户昵称
	Type     int    `json:"password"` // 用户类型
}

// @title             GetInfo
// @description       获取用户信息
// @auth              高宏宇         2022/2/12
// @param             c             请求句柄
func GetInfo(c *gin.Context) {
	userName, _ := c.Cookie("camp-session")

	redisConn := redisPool.Get()
	defer redisConn.Close()
	var dbsearchResult string
	var user *database.User

	// 读取redis或database获取user
	val, err := redisConn.Do("HMGET", hash, userName)
	if err == redis.Nil { // redis查询结果为空
		dbsearchResult, user = database.GetUserInfoByName(userName)
		if dbsearchResult == "Success" {
			redisConn.Do("HSET", hash, userName, user)
		} else {
			c.JSON(types.UserNotExisted, gin.H{"status": types.UserNotExisted})
			return
		}
	} else if err != nil { // redis查询出现错误
		tools.LogMsg(err)
		c.JSON(types.UnknownError, gin.H{"status": types.UnknownError})
		return
	} else { // redis查询到结果
		user = val.(*database.User)
	}

	json := responseJson{
		Id:       user.Id,
		Nickname: user.Nickname,
		Type:     user.Type,
	}
	c.JSON(http.StatusOK, json)
}
