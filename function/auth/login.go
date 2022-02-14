package auth

import (
	"backend/database"
	"backend/function/selectCourse"
	"backend/tools"
	"backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// 请求信息
type requestJson struct {
	Username string `form:"userName" json:"userName" binding:"required"` // 用户名
	Pssword  string `form:"password" json:"password" binding:"required"` // 密码
}

// @title             Login
// @description       登录
// @auth              高宏宇         2022/2/12
// @param             c             请求句柄
func Login(c *gin.Context) {
	// 已经登录无需再次身份验证
	if _, err := c.Cookie("camp-session"); err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "200"})
		return
	}

	var json requestJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redisConn := selectCourse.RedisPool.Get()
	defer redisConn.Close()
	var dbsearchResult string
	var user *database.User

	// 读取redis或database获取user
	val, err := redisConn.Do("HMGET", hash, json.Username)
	if err == redis.Nil { // redis查询结果为空
		dbsearchResult, user = database.GetUserInfoByName(json.Username)
		if dbsearchResult == "Success" {
			redisConn.Do("HSET", hash, json.Username, user)
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

	status := http.StatusOK
	if user.IsValid == 0 {
		status = types.UserHasDeleted
	} else if json.Pssword != user.Password {
		status = types.WrongPassword
	}
	c.JSON(status, gin.H{"status": status})

	if status == http.StatusOK {
		c.SetCookie("camp-session", json.Username, 3600, "/", "localhost", false, true)
	}
}
