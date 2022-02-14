package auth

import (
	"backend/database"
	"backend/function/selectCourse"
	"backend/tools"
	"backend/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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
	userId, err := c.Cookie("camp-session")
	if err != nil {
		whoAmIResponse := types.WhoAmIResponse{
			Code: types.LoginRequired,
			Data: types.TMember{},
		}
		c.JSON(types.LoginRequired, whoAmIResponse)
		return
	}

	redisConn := selectCourse.RedisPool.Get()
	defer redisConn.Close()
	var user *database.User
	// 读取redis或database获取user
	val, err := redisConn.Do("HMGET", hash, userId)
	user = val.(*database.User)
	// redis查询出现错误
	if err != nil {
		tools.LogMsg(err)
		whoAmIResponse := types.WhoAmIResponse{
			Code: types.UnknownError,
			Data: types.TMember{},
		}
		c.JSON(types.UnknownError, whoAmIResponse)
		return
	}
	// redis查询结果为空
	if err == redis.Nil {
		user = database.GetUserInfoById(userId)
		if user == nil {
			whoAmIResponse := types.WhoAmIResponse{
				Code: types.UserNotExisted,
				Data: types.TMember{},
			}
			c.JSON(types.UserNotExisted, whoAmIResponse)
			return
		}
		redisConn.Do("HSET", hash, user.Id, user)
		redisConn.Do("HSET", hash, user.Name, user)
	}

	whoAmIResponse := types.WhoAmIResponse{
		Code: http.StatusOK,
		Data: types.TMember{
			UserID:   strconv.Itoa(user.Id),
			Nickname: user.Nickname,
			Username: user.Name,
			UserType: user.Type,
		},
	}
	c.JSON(http.StatusOK, whoAmIResponse)
}
