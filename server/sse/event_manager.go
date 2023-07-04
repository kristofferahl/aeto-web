package sse

import (
	"log"
	"sync"
)

type ServerSentEvent interface {
	Payload() ([]byte, error)
}

type EventManager struct {
	subscribers map[string]map[chan<- ServerSentEvent]struct{}
	mu          sync.Mutex
}

func NewEventManager() *EventManager {
	return &EventManager{
		subscribers: make(map[string]map[chan<- ServerSentEvent]struct{}),
	}
}

func (em *EventManager) Subscribe(eventType string, ch chan<- ServerSentEvent) {
	log.Println("subscribing to", eventType)
	em.mu.Lock()
	defer em.mu.Unlock()

	subscribers, ok := em.subscribers[eventType]
	if !ok {
		subscribers = make(map[chan<- ServerSentEvent]struct{})
		em.subscribers[eventType] = subscribers
	}

	subscribers[ch] = struct{}{}
}

func (em *EventManager) Unsubscribe(eventType string, ch chan<- ServerSentEvent) {
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

func (em *EventManager) Publish(eventType string, event ServerSentEvent) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if subscribers, ok := em.subscribers[eventType]; ok {
		for ch := range subscribers {
			ch <- event
		}
	}
}
