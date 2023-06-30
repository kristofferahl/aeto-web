package sse

import "sync"

type Event struct {
	Type    string      // Type of the event
	Payload interface{} // Event data payload
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

	if subscribers, ok := em.subscribers[eventType]; ok {
		for ch := range subscribers {
			ch <- event
		}
	}
}
