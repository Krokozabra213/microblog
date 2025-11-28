package logger

import (
	"fmt"
	"log"
	"os"
)

// этот интерфейс выступает как нобор правил
// Логгер на то что поставлен лайк
// Логгер на снятие лайка
// Логгер на ошибку в процессе
type LikeLogger interface {
	LogLike(postID string, userID string)
	LogUnlike(postID string, userID string)
	logerr(postID string, userID string, err error)
}

type consoleLikeLogger struct {
	logger *log.Logger
}

func NewConsoleLikeLogger() *consoleLikeLogger {
	return &consoleLikeLogger{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (l *consoleLikeLogger) logLike(postID string, userID string) {
	l.logger.Printf(fmt.Sprintf("User %s liked post %s", userID, postID))
}

func (l *consoleLikeLogger) logUnlike(postID string, userID string) {
	l.logger.Printf(fmt.Sprintf("User %s unliked post %s", userID, postID))
}

func (l *consoleLikeLogger) logerr(postID string, userID string, err error) {
	l.logger.Printf(fmt.Sprintf("Error while liking post %s for user %s: %s", postID, userID, err))
}
