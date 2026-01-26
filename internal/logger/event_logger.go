package logger

import (
	"log"
	"sync"
	"time"
)

type EventLogger struct {
	channel chan EventMessage
	closed  chan struct{}
	once    sync.Once
	wg      sync.WaitGroup
}

func NewEventLogger(buffer int) *EventLogger {
	logger := &EventLogger{
		channel: make(chan EventMessage, buffer),
		closed:  make(chan struct{}),
	}

	logger.wg.Add(1)
	go func() {
		defer logger.wg.Done()
		for event := range logger.channel {
			log.Println(event.EventMessage())
		}
	}()

	return logger
}

func (l *EventLogger) Log(event EventMessage) {

	select {
	case <-l.closed:
		return
	default:
	}

	defer func() {
		recover()
	}()

	select {
	case l.channel <- event:
	default:
		log.Printf("WARNING: Event channel full, event dropped. Channel size: %d/%d",
			len(l.channel), cap(l.channel))
	}
}

func (l *EventLogger) GracefullShutdown(timeout time.Duration) {
	l.once.Do(func() {
		close(l.closed)  // Сигнал для Log() - больше не принимаем
		close(l.channel) // Сигнал для горутины - завершайся (после drain)
	})

	// Ждём завершения с таймаутом
	done := make(chan struct{})
	go func() {
		l.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Logger: graceful shutdown complete")
	case <-time.After(timeout):
		log.Println("Logger: shutdown timeout, some events may be lost")
	}
}
