package logger

import (
	"context"
	"log"
	"os"
)

type Options struct {
	AppEnv  string
	AppName string
}

var logger *log.Logger

func New(options Options) {
	logger = log.New(os.Stdout, "["+options.AppName+"] ", log.LstdFlags|log.Lshortfile)
}

func Info(ctx context.Context, message string) {
	if logger == nil {
		log.Println("[INFO]", message)
		return
	}
	logger.Println("[INFO]", message)
}

func Error(ctx context.Context, message string) {
	if logger == nil {
		log.Println("[ERROR]", message)
		return
	}
	logger.Println("[ERROR]", message)
}

func Debug(ctx context.Context, message string) {
	if logger == nil {
		log.Println("[DEBUG]", message)
		return
	}
	logger.Println("[DEBUG]", message)
}

func Warn(ctx context.Context, message string) {
	if logger == nil {
		log.Println("[WARN]", message)
		return
	}
	logger.Println("[WARN]", message)
}

