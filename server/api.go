package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/teacat/jsonfilter"
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
		r.Get("/tenants/{namespace}/{name}", getTenant(client))
		r.Get("/resourcesets", listResourceSets(client))
		r.Get("/resourcesets/{namespace}/{name}", getResourceSet(client))
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

		tl, err := client.CoreV1Alpha1("aeto").ListTenants(v1.ListOptions{})
		if hasErr(w, err) {
			return
		}

		data, err := json.Marshal(tl)
		if hasErr(w, err) {
			return
		}

		res, err := jsonfilter.Filter(data, "items(metadata(annotations,creationTimestamp,finalizers,generation,name,namespace,resourceVersion,uid),spec,status)")
		if hasErr(w, err) {
			return
		}

		w.Write(res)
	}
}

func getTenant(client *AetoClient) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		namespace := chi.URLParam(req, "namespace")
		name := chi.URLParam(req, "name")

		if namespace != "" && name != "" {
			rs, err := client.CoreV1Alpha1(namespace).GetTenant(name)
			if hasErr(w, err) {
				// TODO: Handle 404
				return
			}

			data, err := json.Marshal(rs)
			if hasErr(w, err) {
				return
			}

			w.Write(data)
		} else {
			w.WriteHeader(400)
			w.Write([]byte(""))
		}
	}
}

func listResourceSets(client *AetoClient) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tl, err := client.CoreV1Alpha1("aeto").ListResourceSets(v1.ListOptions{})
		if hasErr(w, err) {
			return
		}

		data, err := json.Marshal(tl)
		if hasErr(w, err) {
			return
		}

		res, err := jsonfilter.Filter(data, "items(metadata(annotations,creationTimestamp,finalizers,generation,name,namespace,resourceVersion,uid),spec,status)")
		if hasErr(w, err) {
			return
		}

		w.Write(res)
	}
}

func getResourceSet(client *AetoClient) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		namespace := chi.URLParam(req, "namespace")
		name := chi.URLParam(req, "name")

		if namespace != "" && name != "" {
			rs, err := client.CoreV1Alpha1(namespace).GetResourceSet(name)
			if hasErr(w, err) {
				// TODO: Handle 404
				return
			}

			data, err := json.Marshal(rs)
			if hasErr(w, err) {
				return
			}

			w.Write(data)
		} else {
			w.WriteHeader(400)
			w.Write([]byte(""))
		}
	}
}
