package logger

import (
	"fmt"
	"log"
	"os"
)

type RegistLog interface {
	LogRegist(UserID string)
	LogRegistError(UserID string, err error)
}

type consoleRegistLogger struct {
	logger *log.Logger
}

func NewconsoleRegistLogger() RegistLog {
	return &consoleRegistLogger{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (l *consoleRegistLogger) LogRegist(UserID string) {
	l.logger.Println(fmt.Sprintf("[REGISTRATION] User %s registered", UserID))
}

func (l *consoleRegistLogger) LogRegistError(UserID string, err error) {
	l.logger.Println(fmt.Sprintf("[REGISTRATION ERROR] User %s error: %v", UserID, err))
}
