package sse

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func HandleSSE(em *EventManager, w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	cid := r.URL.Query().Get("cid")
	if len(cid) == 0 {
		cid = string(uuid.New().String())[0:8]
	}

	log.Println("/sse, establishing connection", cid)

	ch := make(chan ServerSentEvent)
	em.Subscribe("change", ch)

	defer func() {
		log.Println("/sse, exiting handler, connection closed", cid)
		em.Unsubscribe("change", ch)
		close(ch)
	}()

	for {
		select {
		case event := <-ch:
			log.Println("sse, write event to stream", cid)
			p, err := event.Payload()
			if err != nil {
				log.Println("sse, error getting the payload of the event", cid)
				continue
			}

			eventData := fmt.Sprintf("data: %s\n\n", string(p))

			_, err = w.Write([]byte(eventData))
			if err != nil {
				log.Println("sse, error writing to http.ResponseWriter", cid)
				continue
			}

			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			log.Println("http request closed", cid)
			return
		}
	}
}
