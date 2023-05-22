package server

import (
	rest "k8s.io/client-go/rest"
)

type AetoClient struct {
	restConfig    *rest.Config
	corev1Alpha1  *rest.RESTClient
	eventv1Alpha1 *rest.RESTClient
}

func NewForConfig(c *rest.Config) (*AetoClient, error) {
	client := &AetoClient{
		restConfig: c,
	}

	corev1Alpha1Client, err := client.NewCoreV1Alpha1Client()
	if err != nil {
		return nil, err
	}

	eventv1Alpha1Client, err := client.NewEventV1Alpha1Client()
	if err != nil {
		return nil, err
	}

	return &AetoClient{
		restConfig:    c,
		corev1Alpha1:  corev1Alpha1Client,
		eventv1Alpha1: eventv1Alpha1Client,
	}, nil
}
