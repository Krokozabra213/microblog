package queue

import (
	"errors"
	"microblog/internal/logger"
	"sync"
	"time"
)

var ErrShutdownTimeout = errors.New("shutdown timeout")

type PostService interface {
	LikePost(postID, userID string) error
}

type Logger interface {
	Log(event logger.EventMessage)
}

type LikeEvent struct {
	PostID string
	UserID string
}

type LikeQueue struct {
	channel     chan LikeEvent
	postService PostService
	eventLogger Logger

	closed    chan struct{}
	closeOnce sync.Once
	wg        sync.WaitGroup
}

func NewLikeQueue(eventLogger Logger, postService PostService, eventBuffer int) *LikeQueue {
	ch := make(chan LikeEvent, eventBuffer)

	q := &LikeQueue{
		channel:     ch,
		postService: postService,
		eventLogger: eventLogger,
		closed:      make(chan struct{}),
	}

	q.wg.Add(1)
	go q.worker()
	return q
}

func (q *LikeQueue) worker() {
	defer q.wg.Done()
	for event := range q.channel {
		q.processEvent(event)
	}
}

func (q *LikeQueue) processEvent(e LikeEvent) {
	err := q.postService.LikePost(e.PostID, e.UserID)

	eventType := logger.PostLiked
	message := logger.PostLikedMessage

	if err != nil {
		eventType = logger.PostLikedErr
		message = err.Error()
	}

	q.eventLogger.Log(logger.EventPost{
		Type:      eventType,
		AuthorID:  e.UserID,
		PostID:    e.PostID,
		Message:   message,
		Timestamp: time.Now(),
	})
}

func (q *LikeQueue) Enqueue(event LikeEvent) bool {
	// Быстрая проверка
	select {
	case <-q.closed:
		return false
	default:
	}

	// Защита от panic
	defer func() {
		recover()
	}()

	select {
	case q.channel <- event:
		return true
	case <-q.closed:
		return false
	}
}

func (q *LikeQueue) GracefullShutdown(timeout time.Duration) error {
	q.closeOnce.Do(func() {
		close(q.closed)
		close(q.channel)
	})

	// Ждём завершения с таймаутом
	done := make(chan struct{})
	go func() {
		q.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return ErrShutdownTimeout
	}
}
