package logger

import (
	"log"
	"os"
)

// Logger for default logger
var logger *log.Logger

const LogFile = "text.log"

// InitLogger initializes logger
func InitLogger() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	f, err := os.OpenFile(LogFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	// defer f.Close()

	if os.Getenv("STD_LOG") != "" {
		logger = log.New(os.Stdout, "PRIME ", log.LstdFlags)
		return
	}

	logger = log.New(f, "PRIME ", log.LstdFlags)
}

// Println prints logs
func Println(v ...interface{}) {
	logger.Println(v...)
}
