package tools

import (
	"log"
	"os"
)

var Logger *log.Logger

func Init() {
	logFile, err := os.OpenFile("./app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// fmt.Println("failed to open file!")
		return
	}
	Logger = log.New(logFile, "[DEBUG]", log.Ldate|log.Lshortfile)
}

func LogMsg(arg ...interface{}) {
	Logger.Println(arg...)
}

func LogfMsg(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}
