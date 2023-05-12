package server

import (
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func addUiRoutes(s *Server, router *chi.Mux) {
	router.Handle("/*", handleUI(s))
}

func handleUI(s *Server) http.HandlerFunc {
	assetFS := s.getAssets()
	assetHandler := http.FileServer(http.FS(assetFS))
	redirector := createRedirector(assetFS)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		extension := filepath.Ext(req.URL.Path)
		// We use the golang http.FileServer for static file requests.
		// This will return a 404 on normal page requests, ie /kustomizations and /sources.
		// Redirect all non-file requests to index.html, where the JS routing will take over.
		if extension == "" {
			redirector(w, req)
			return
		}
		assetHandler.ServeHTTP(w, req)
	})
}

func createRedirector(fsys fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexPage, err := fsys.Open("index.html")

		if err != nil {
			log.Printf("could not open index.html page, %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		stat, err := indexPage.Stat()
		if err != nil {
			log.Printf("could not get index.html stat, %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bt := make([]byte, stat.Size())
		_, err = indexPage.Read(bt)

		if err != nil {
			log.Printf("could not read index.html, %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(bt)

		if err != nil {
			log.Printf("error writing index.html, %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
