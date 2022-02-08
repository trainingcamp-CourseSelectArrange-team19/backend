package main

import (
	"backend/database"	
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	r := gin.Default()
	RegisterRouter(r)
	r.Run(":8000")
	
}
