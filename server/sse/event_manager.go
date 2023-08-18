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

func (em *EventManager) Subscribe(eventstream string, ch chan<- ServerSentEvent) {
	log.Println("subscribing to", eventstream)
	em.mu.Lock()
	defer em.mu.Unlock()

	subscribers, ok := em.subscribers[eventstream]
	if !ok {
		subscribers = make(map[chan<- ServerSentEvent]struct{})
		em.subscribers[eventstream] = subscribers
	}

	subscribers[ch] = struct{}{}
}

func (em *EventManager) Unsubscribe(eventstream string, ch chan<- ServerSentEvent) {
	log.Println("unsubscribing from", eventstream)
	em.mu.Lock()
	defer em.mu.Unlock()

	if subscribers, ok := em.subscribers[eventstream]; ok {
		delete(subscribers, ch)
		if len(subscribers) == 0 {
			delete(em.subscribers, eventstream)
		}
	}
}

func (em *EventManager) Publish(eventstream string, event ServerSentEvent) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if subscribers, ok := em.subscribers[eventstream]; ok {
		for ch := range subscribers {
			ch <- event
		}
	}
}
