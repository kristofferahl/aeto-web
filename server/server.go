package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/kristofferahl/aeto-web/server/sse"
)

var eventManager *sse.EventManager

type Server struct {
	EmbeddedFiles     embed.FS
	EmbeddedFilesPath string
	ClusterConfig     bool
}

func (s *Server) Run() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))

	addUiRoutes(s, r)
	addApiRoutes(s, r, eventManager)

	log.Println("aeto server is listening on port 9000...")
	err := http.ListenAndServe(":9000", r)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) getAssets() fs.FS {
	f, err := fs.Sub(s.EmbeddedFiles, s.EmbeddedFilesPath)

	if err != nil {
		panic(err)
	}
	return f
}

func publishKeepAlive() {
	eventManager.Publish("change", &ApiEvent{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Type:      "KeepAlive",
	})
	time.Sleep(15 * time.Second)
	go publishKeepAlive()
}

func init() {
	eventManager = sse.NewEventManager()
	go publishKeepAlive()
}
