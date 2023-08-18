package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/kristofferahl/aeto-web/server/sse"
)

type ApiEvent struct {
	Timestamp string `json:"ts"`
	Type      string `json:"type"`
}

func (e *ApiEvent) Payload() ([]byte, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return b, err
}

func addApiRoutes(s *Server, router *chi.Mux, em *sse.EventManager) {
	restConfig, err := getRestConfig(s.ClusterConfig)
	if err != nil {
		panic(err)
	}

	client, err := NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	// TODO: Make this configurable
	operatorNamespace := "aeto"

	WatchKubernetesEvents(client.clientset, operatorNamespace)

	if err := client.AcmAwsV1Alpha1(operatorNamespace).Watch(); err != nil {
		panic(err)
	}
	if err := client.CoreV1Alpha1(operatorNamespace).Watch(); err != nil {
		panic(err)
	}
	if err := client.EventV1Alpha1(operatorNamespace).Watch(); err != nil {
		panic(err)
	}
	if err := client.Route53AwsV1Alpha1(operatorNamespace).Watch(); err != nil {
		panic(err)
	}
	if err := client.SustainabilityV1Alpha1(operatorNamespace).Watch(); err != nil {
		panic(err)
	}

	router.Group(func(r chi.Router) {
		// r.Use(middleware.Timeout(10 * time.Second))

		r.Get("/api/sse", func(w http.ResponseWriter, r *http.Request) {
			sse.HandleSSE(em, w, r)
		})
	})

	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.Timeout(2 * time.Second))

		r.Get("/dashboard", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			tenants, err := client.CoreV1Alpha1(operatorNamespace).ListTenants()
			if hasErr(w, err) {
				return
			}
			blueprints, err := client.CoreV1Alpha1(operatorNamespace).ListBlueprints()
			if hasErr(w, err) {
				return
			}
			certificates, err := client.AcmAwsV1Alpha1(operatorNamespace).ListCertificates()
			if hasErr(w, err) {
				return
			}
			hostedzones, err := client.Route53AwsV1Alpha1(operatorNamespace).ListHostedZones()
			if hasErr(w, err) {
				return
			}
			savingspolicies, err := client.SustainabilityV1Alpha1(operatorNamespace).ListSavingsPolicies()
			if hasErr(w, err) {
				return
			}

			dashboard := struct {
				Tenants         int `json:"tenants"`
				Blueprints      int `json:"blueprints"`
				Certificates    int `json:"certificates"`
				HostedZones     int `json:"hostedzones"`
				SavingsPolicies int `json:"savingspolicies"`
			}{
				Tenants:         len(tenants.Items),
				Blueprints:      len(blueprints.Items),
				Certificates:    len(certificates.Items),
				HostedZones:     len(hostedzones.Items),
				SavingsPolicies: len(savingspolicies.Items),
			}

			data, err := json.Marshal(dashboard)
			if hasErr(w, err) {
				return
			}

			w.Write(data)
		})

		r.Get("/tenants", listResource(client, func() (interface{}, error) {
			return client.CoreV1Alpha1(operatorNamespace).ListTenants()
		}))
		r.Get("/tenants/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.CoreV1Alpha1(namespace).GetTenant(name)
		}))

		r.Get("/blueprints", listResource(client, func() (interface{}, error) {
			return client.CoreV1Alpha1(operatorNamespace).ListBlueprints()
		}))
		r.Get("/blueprints/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.CoreV1Alpha1(namespace).GetBlueprint(name)
		}))

		r.Get("/resourcesets", listResource(client, func() (interface{}, error) {
			return client.CoreV1Alpha1(operatorNamespace).ListResourceSets()
		}))
		r.Get("/resourcesets/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.CoreV1Alpha1(namespace).GetResourceSet(name)
		}))

		r.Get("/resourcetemplates", listResource(client, func() (interface{}, error) {
			return client.CoreV1Alpha1(operatorNamespace).ListResourceTemplates()
		}))
		r.Get("/resourcetemplates/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.CoreV1Alpha1(namespace).GetResourceTemplate(name)
		}))

		r.Get("/eventstreamchunks", listResource(client, func() (interface{}, error) {
			return client.EventV1Alpha1(operatorNamespace).ListEventStreamChunks()
		}))
		r.Get("/eventstreamchunks/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.EventV1Alpha1(namespace).GetEventStreamChunk(name)
		}))

		r.Get("/savingspolicies", listResource(client, func() (interface{}, error) {
			return client.SustainabilityV1Alpha1(operatorNamespace).ListSavingsPolicies()
		}))
		r.Get("/savingspolicies/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.SustainabilityV1Alpha1(namespace).GetSavingsPolicy(name)
		}))

		r.Get("/certificates", listResource(client, func() (interface{}, error) {
			return client.AcmAwsV1Alpha1(operatorNamespace).ListCertificates()
		}))
		r.Get("/certificates/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.AcmAwsV1Alpha1(namespace).GetCertificate(name)
		}))
		r.Get("/certificateconnectors", listResource(client, func() (interface{}, error) {
			return client.AcmAwsV1Alpha1(operatorNamespace).ListCertificateConnectors()
		}))
		r.Get("/certificateconnectors/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.AcmAwsV1Alpha1(namespace).GetCertificateConnector(name)
		}))

		r.Get("/hostedzones", listResource(client, func() (interface{}, error) {
			return client.Route53AwsV1Alpha1(operatorNamespace).ListHostedZones()
		}))
		r.Get("/hostedzones/{namespace}/{name}", getResource(client, func(namespace, name string) (interface{}, error) {
			return client.Route53AwsV1Alpha1(namespace).GetHostedZone(name)
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

		res, err := ApplyResourceListFilter(data)
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

			res, err := ApplyResourceFilter("%s", data)
			if hasErr(w, err) {
				return
			}

			w.Write(res)
		} else {
			w.WriteHeader(400)
			w.Write([]byte(""))
		}
	}
}
