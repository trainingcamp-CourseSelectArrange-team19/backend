package main

import (
	"backend/database"
	"backend/function/selectCourse"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	selectCourse.InitRedisConfig()
	r := gin.Default()
	pprof.Register(r)
	RegisterRouter(r)
	r.Run(":8000")
	
}
