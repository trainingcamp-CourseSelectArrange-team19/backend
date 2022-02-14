package auth

import (
	"github.com/garyburd/redigo/redis"
)

var (
	redisPool *redis.Pool
	hash      string = "users"
)

// @title             NewPool
// @description       初始化redis连接池
// @auth              高宏宇         2022/2/13
// @return			  *redis.Pool	redis连接池
func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   10000, // 最大空闲连接数
		MaxActive: 12000, // 最大连接数
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// @title             InitRedisConfig
// @description       初始化redis
// @auth              高宏宇         2022/2/13
func InitRedisConfig() {
	redisPool = NewPool()
}
