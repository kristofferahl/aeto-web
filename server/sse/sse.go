package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandleSSE(em *EventManager, w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	log.Println("/sse, establishing connection")

	ch := make(chan Event)
	em.Subscribe("change", ch)

	defer func() {
		log.Println("/sse, connection closing")
		em.Unsubscribe("change", ch)
		close(ch)
	}()

	for {
		select {
		case event := <-ch:
			// Send the event to the SSE client
			b, err := json.Marshal(event.Payload)
			if err != nil {
				log.Println("/sse, error writing doing marshalling")
				continue
			}
			eventData := fmt.Sprintf("data: %s\n\n", string(b))

			_, err = w.Write([]byte(eventData))
			if err != nil {
				// Handle error
				log.Println("/sse, error writing to responseWriter")
				continue
			}

			// if f, ok := w.(http.Flusher); ok {
			// 	f.Flush()
			// }
		case <-r.Context().Done():
			log.Println("Context done")
			return
		}
	}
}
