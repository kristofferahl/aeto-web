package server

import (
	"os"
	"path/filepath"

	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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
