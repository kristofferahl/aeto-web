package sse

import (
	"log"
	"sync"
	"time"
)

type Event struct {
	Type      string      `json:"type"`    // Type of the event
	Payload   interface{} `json:"payload"` // Event data payload
	Timestamp string      `json:"ts"`
}

type EventManager struct {
	subscribers map[string]map[chan<- Event]struct{}
	mu          sync.Mutex
}

func NewEventManager() *EventManager {
	return &EventManager{
		subscribers: make(map[string]map[chan<- Event]struct{}),
	}
}

func (em *EventManager) Subscribe(eventType string, ch chan<- Event) {
	log.Println("subscribing to", eventType)
	em.mu.Lock()
	defer em.mu.Unlock()

	subscribers, ok := em.subscribers[eventType]
	if !ok {
		subscribers = make(map[chan<- Event]struct{})
		em.subscribers[eventType] = subscribers
	}

	subscribers[ch] = struct{}{}
}

func (em *EventManager) Unsubscribe(eventType string, ch chan<- Event) {
	log.Println("unsubscribing from", eventType)
	em.mu.Lock()
	defer em.mu.Unlock()

	if subscribers, ok := em.subscribers[eventType]; ok {
		delete(subscribers, ch)
		if len(subscribers) == 0 {
			delete(em.subscribers, eventType)
		}
	}
}

func (em *EventManager) Publish(eventType string, event Event) {
	em.mu.Lock()
	defer em.mu.Unlock()
	event.Timestamp = time.Now().UTC().Format(time.RFC3339)

	if subscribers, ok := em.subscribers[eventType]; ok {
		for ch := range subscribers {
			ch <- event
		}
	}
}
