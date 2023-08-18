package server

import (
	"encoding/json"
	"os"
	"path/filepath"

	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesEvent struct {
	Timestamp string        `json:"ts"`
	EventType string        `json:"type"`
	Reason    string        `json:"reason"`
	Message   string        `json:"message"`
	Resource  EventResource `json:"resource"`
}

type EventResource struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
}

func (e *KubernetesEvent) Payload() ([]byte, error) {
	b, err := json.Marshal(e)
	return b, err
}

func getRestConfig(inClusterConfig bool) (*rest.Config, error) {
	var err error
	var restConfig *rest.Config

	if inClusterConfig {
		restConfig, err = rest.InClusterConfig()
	} else {
		kubeconfig := filepath.Join(homeDir(), ".kube", "config")
		if path := os.Getenv("KUBECONFIG"); path != "" {
			kubeconfig = path
		}
		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return restConfig, err
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
