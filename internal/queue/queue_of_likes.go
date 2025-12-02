package queue

import (
	"microblog/internal/logger"
	"microblog/internal/service"
)

type LikeQueue struct {
	channel     chan LikeEvent
	postService *service.PostService
	logger      logger.LikeLogger
}
type LikeEvent struct {
	PostID string
	UserID string
}

func (q *LikeQueue) Start() {

	for {
		LikeEvent := <-q.channel

		_, err := q.postService.LikePost(LikeEvent.PostID, LikeEvent.UserID)

		if err != nil {
			q.logger.LogError(LikeEvent.PostID, LikeEvent.UserID, err)
		} else {
			q.logger.LogLike(LikeEvent.PostID, LikeEvent.UserID)
		}
	}

}
func NewLikeQueue(postService *service.PostService, logger logger.LikeLogger) *LikeQueue {
	ch := make(chan LikeEvent, 100)

	q := &LikeQueue{
		channel:     ch,
		postService: postService,
		logger:      logger,
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
