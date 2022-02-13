package main

import (
	"backend/auth"
	"backend/database"
	"backend/function/selectCourse"
	"reflect"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	selectCourse.InitRedisConfig()
	auth.InitRedisConfig()
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			<-ticker.C
			for val := range selectCourse.SuccessSet.Iter() {
				selectCourse.InsertSchedule(reflect.ValueOf(val).String())
			}
			selectCourse.SuccessSet.Clear()
		}
	}()
	defer ticker.Stop()
	r := gin.Default()
	pprof.Register(r)
	RegisterRouter(r)
	r.Run(":8000")
}
