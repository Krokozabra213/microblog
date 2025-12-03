package logger

import (
	"log"
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
	l.channel <- event
}

func (l *EventLogger) Close() {

	close(l.channel)
}
