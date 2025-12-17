package queue

import (
	"microblog/internal/logger"
	"microblog/internal/service"
	"time"
)

type LikeQueue struct {
	channel     chan LikeEvent
	postService *service.PostService
	eventLogger *logger.EventLogger
}

type LikeEvent struct {
	PostID string
	UserID string
}

func (q *LikeQueue) Start() {
	for e := range q.channel {
		_, err := q.postService.LikePost(e.PostID, e.UserID)

		if err != nil {
			event := logger.Event{
				Type:      "LIKE_ERROR",
				UserID:    e.UserID,
				PostID:    e.PostID,
				Message:   err.Error(),
				Timestemp: time.Now(),
			}
			q.eventLogger.Log(event)
		} else {
			event := logger.Event{
				Type:      "LIKE",
				UserID:    e.UserID,
				PostID:    e.PostID,
				Message:   "Post liked successfully",
				Timestemp: time.Now(),
			}
			q.eventLogger.Log(event)
		}
	}
}

func NewLikeQueue(postService *service.PostService, eventLogger *logger.EventLogger) *LikeQueue {
	ch := make(chan LikeEvent, 100)

	q := &LikeQueue{
		channel:     ch,
		postService: postService,
		eventLogger: eventLogger,
	}

	go q.Start()
	return q
}

func (q *LikeQueue) Enqueue(event LikeEvent) {
	q.channel <- event
}

func (q *LikeQueue) Close() {
	close(q.channel)
}
