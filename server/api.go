package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func addApiRoutes(s *Server, router *chi.Mux) {
	restConfig, err := getRestConfig(s.InClusterConfig)
	if err != nil {
		panic(err)
	}

	client, err := NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.Timeout(60 * time.Second))

		r.Get("/tenants", listTenants(client))
	})
}

func hasErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println("an unhandled error occured,", err)
		w.WriteHeader(500)
		return true
	}
	return false
}

func listTenants(client *AetoClient) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tl, err := client.CoreV1Alpha1("aeto").List(v1.ListOptions{})
		if hasErr(w, err) {
			return
		}

		res, err := json.Marshal(tl)
		if hasErr(w, err) {
			return
		}

		w.Write(res)
	}
}
