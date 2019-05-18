package utils

import (
	"log"
	"os"
	"strings"
)

var (
	logger *log.Logger
	env    string
)

func init() {
	env = strings.ToLower(os.Getenv("env"))
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

func Debug(format string, v ...interface{}) {
	if env == "development" {
		logger.Printf("[DEBUG] "+format, v...)
	}
}

func Info(format string, v ...interface{}) {
	logger.Printf("[INFO] "+format, v...)
}

func Warning(format string, v ...interface{}) {
	logger.Printf("[WARN] "+format, v...)
}

func Error(format string, v ...interface{}) {
	logger.Fatalf("[ERROR] "+format, v...)
	os.Exit(1)
}
