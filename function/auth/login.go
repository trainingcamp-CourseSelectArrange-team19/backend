package auth

import (
	"backend/database"
	"backend/tools"
	"backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type requestJson struct {
	Username string `form:"userName" json:"userName" binding:"required"`
	Pssword  string `form:"password" json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	if _, err := c.Cookie("camp-session"); err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "200"})
		return
	}

	var json requestJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redisConn := redisPool.Get()
	defer redisConn.Close()
	var dbsearchResult string
	var user *database.User

	val, err := redisConn.Do("HMGET", hash, json.Username)
	if err == redis.Nil {
		dbsearchResult, user = database.GetUserInfoByName(json.Username)
		if dbsearchResult == "Success" {
			redisConn.Do("HSET", hash, json.Username, user)
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
