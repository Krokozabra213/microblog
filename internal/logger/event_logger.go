package logger

import (
	"fmt"
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
	mu.Lock()
	select {
	case l.channel <- event:
		fmt.Println("Event logged")
	default:
		log.Printf("Предупреждение: не удалось залогировать событие, канал переполнен: %+v", event)

	}
	mu.Unlock()
}

func (l *EventLogger) Close() {

	close(l.channel)
}
