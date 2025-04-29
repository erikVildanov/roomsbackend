package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Grey   = "\033[90m"
	Green  = "\033[32m"
)

// Infof выводит информационное сообщение.
func Infof(format string, args ...any) {
	print("[INFO]", Green, format, args...)
}

// Warnf выводит предупреждение.
func Warnf(format string, args ...any) {
	print("[WARN]", Yellow, format, args...)
}

// Errorf выводит ошибку.
func Errorf(format string, args ...any) {
	print("[ERROR]", Red, format, args...)
}

// Debugf выводит отладочную информацию.
func Debugf(format string, args ...any) {
	if os.Getenv("DEBUG") == "true" {
		print("[DEBUG]", Cyan, format, args...)
	}
}

// print — базовый форматтер лога.
func print(prefix, color, format string, args ...any) {
	timestamp := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, args...)
	log.Printf("%s %s%s%s %s", Grey+timestamp+Reset, color, prefix, Reset, msg)
}
