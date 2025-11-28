package logger

import (
	"fmt"
	"log"
	"os"
)

type LikeLogger interface {
	LogLike(postID string, userID string)
	LogUnlike(postID string, userID string)
	LogError(postID string, userID string, err error)
}

type consoleLikeLogger struct {
	logger *log.Logger
}

func NewConsoleLikeLogger() LikeLogger {
	return &consoleLikeLogger{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (l *consoleLikeLogger) LogLike(postID string, userID string) {
	l.logger.Printf(fmt.Sprintf("[LIKE] User %s liked post %s", userID, postID))
}

func (l *consoleLikeLogger) LogUnlike(postID string, userID string) {
	l.logger.Printf(fmt.Sprintf("[UNLIKE] User %s unliked post %s", userID, postID))
}

func (l *consoleLikeLogger) LogError(postID string, userID string, err error) {
	l.logger.Printf(fmt.Sprintf("[ERROR] User %s post %s error: %v", userID, postID, err))
}
