package logger

import (
	"log"
	"sync"
	"time"
)

type Event struct {
	Type      string
	UserID    string
	PostID    string
	Message   string
	Timestemp time.Time
}

type EventLogger struct {
	channel chan Event
}

var mu sync.Mutex

func NewEventLogger() *EventLogger {
	i := &EventLogger{
		channel: make(chan Event, 100),
	}
	go func() {
		for {
			event := <-i.channel
			log.Print(event)
		}
	}()

	return i
}

func (l *EventLogger) Log(event Event) {

	select {
	case l.channel <- event:
		log.Println("The event was successfully logged.", event)
	default:
		log.Printf("Warning: Failed to log event, channel full: %+v", event)

	}
}

func (l *EventLogger) Close() {
	close(l.channel)
}
