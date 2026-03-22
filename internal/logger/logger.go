package logger

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm/logger"
)

func MakeLoggerFile(path string) *os.File {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("Ошибка создания директории для логов: %v\n", err)
	}
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Ошибка создания лог-файла:", err)
	}

	return logFile
}
func GetAPILogger(path string) *log.Logger {
	file := MakeLoggerFile(path)

	return log.New(file, "", log.LstdFlags)
}

func GetDBLogger(path string) logger.Interface {
	file := MakeLoggerFile(path)

	return logger.New(
		log.New(file, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			// Colorful:                  false,
			// IgnoreRecordNotFoundError: false,
		},
	)
}
