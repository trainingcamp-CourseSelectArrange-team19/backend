package main

import (
/* 	"fmt"
	"runtime/pprof"
	"time"  */
/*  	"os" */

/*  	"backend/database" */
	"github.com/gin-gonic/gin"
)

func main() {
	//tools.Init()
/* 	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			<-ticker.C
			database.CommitAll()
		}
	}()

	defer ticker.Stop()
	database.RedisInitClient()
	database.OpenMysql()
	defer database.RedisClose()
	defer database.CloseMysql()
	database.MysqlCreateTable()
	lst, _ := database.QueryAllRules()
	for index := range *lst {
		fmt.Println((*lst)[index]["id"])
		database.RedisTouchRule((*lst)[index]["id"])
	}
	f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile() */
	r := gin.Default()
	RegisterRouter(r)
	r.Run(":8000")
	
/* 	go r.Run(":8080")
	r2 := gin.Default()
	if os.Getenv("IS_DOCKER") == "1" {
		r2.LoadHTMLGlob("/root/public/index.html")
		database.MysqlCreateTable()
	} else {
		r2.LoadHTMLGlob("./public/index.html")
	}
	//adminRouter(r2)
	r2.Run(":11451") */

	// time.Sleep(50 * time.Second)
	// panic("qaq")

	//r2.Run(":11451")
}
