package logger

import (
	"fmt"
	"log"
	"os"
)

type PostLog interface {
	LogPost(UserID string)
	LogDeletePost(UserID string)
	LogLikePost(UserID string)
}

type consolePostLogger struct {
	logger *log.Logger
}

func NewConsolePostLogger() PostLog {
	return &consolePostLogger{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (l *consolePostLogger) LogPost(UserID string) {
	l.logger.Printf(fmt.Sprintf("[POST] User %s created post", UserID))
}

func (l *consolePostLogger) LogDeletePost(UserID string) {
	l.logger.Printf(fmt.Sprintf("[DELETE POST] User %s deleted post", UserID))
}

func (l *consolePostLogger) LogLikePost(UserID string) {
	l.logger.Printf(fmt.Sprintf("[LIKE POST] User %s liked post", UserID))
}
