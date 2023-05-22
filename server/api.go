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

		r.Get("/tenants", listResource(client, func() (interface{}, error) {
			return client.CoreV1Alpha1("aeto").ListTenants(v1.ListOptions{})
		}))
		r.Get("/tenants/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.CoreV1Alpha1(namespace).GetTenant(name)
		}))

		r.Get("/blueprints", listResource(client, func() (interface{}, error) {
			return client.CoreV1Alpha1("aeto").ListBlueprints(v1.ListOptions{})
		}))
		r.Get("/blueprints/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.CoreV1Alpha1(namespace).GetBlueprint(name)
		}))

		r.Get("/resourcesets", listResource(client, func() (interface{}, error) {
			return client.CoreV1Alpha1("aeto").ListResourceSets(v1.ListOptions{})
		}))
		r.Get("/resourcesets/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.CoreV1Alpha1(namespace).GetResourceSet(name)
		}))

		r.Get("/resourcetemplates", listResource(client, func() (interface{}, error) {
			return client.CoreV1Alpha1("aeto").ListResourceTemplates(v1.ListOptions{})
		}))
		r.Get("/resourcetemplates/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.CoreV1Alpha1(namespace).GetResourceTemplate(name)
		}))

		r.Get("/eventstreamchunks", listResource(client, func() (interface{}, error) {
			return client.EventV1Alpha1("aeto").ListEventStreamChunks(v1.ListOptions{})
		}))
		r.Get("/eventstreamchunks/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.EventV1Alpha1(namespace).GetEventStreamChunk(name)
		}))
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

func listResource(client *AetoClient, list func() (interface{}, error)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		rl, err := list()
		if hasErr(w, err) {
			return
		}

		data, err := json.Marshal(rl)
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

func getResource(client *AetoClient, get func(namespace, name string) (interface{}, error)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		namespace := chi.URLParam(req, "namespace")
		name := chi.URLParam(req, "name")

		if namespace != "" && name != "" {
			rs, err := get(namespace, name)
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
