package auth

import (
	"github.com/garyburd/redigo/redis"
)

var (
	redisPool *redis.Pool
	hash      string = "users"
)

func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   10000,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func InitRedisConfig() {
	redisPool = NewPool()
}
