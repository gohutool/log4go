package main

import (
	"log4go"
	"time"
)

func main() {
	log := log4go.NewLogger()
	defer log.Close()
	log.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())
	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))
}
