package server

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

type AetoClient struct {
	restConfig             *rest.Config
	clientset              *kubernetes.Clientset
	corev1Alpha1           *KubernetesClient
	eventv1Alpha1          *KubernetesClient
	sustainabilityv1Alpha1 *KubernetesClient
	acmAwsV1Alpha1         *KubernetesClient
	route53AwsV1Alpha1     *KubernetesClient
}

func NewForConfig(c *rest.Config) (*AetoClient, error) {
	client := &AetoClient{
		restConfig: c,
	}

	clientset, err := kubernetes.NewForConfig(client.restConfig)
	if err != nil {
		return nil, err
	}

	corev1Alpha1Client, err := client.NewCoreV1Alpha1Client()
	if err != nil {
		return nil, err
	}

	eventv1Alpha1Client, err := client.NewEventV1Alpha1Client()
	if err != nil {
		return nil, err
	}

	sustainabilityv1Alpha1Client, err := client.NewSustainabilityV1Alpha1Client()
	if err != nil {
		return nil, err
	}

	acmAwsV1Alpha1Client, err := client.NewAcmAwsV1Alpha1Client()
	if err != nil {
		return nil, err
	}

	route53AwsV1Alpha1Client, err := client.NewRoute53AwsV1Alpha1Client()
	if err != nil {
		return nil, err
	}

	return &AetoClient{
		restConfig:             c,
		clientset:              clientset,
		corev1Alpha1:           corev1Alpha1Client,
		eventv1Alpha1:          eventv1Alpha1Client,
		sustainabilityv1Alpha1: sustainabilityv1Alpha1Client,
		acmAwsV1Alpha1:         acmAwsV1Alpha1Client,
		route53AwsV1Alpha1:     route53AwsV1Alpha1Client,
	}, nil
}

type KubernetesClient struct {
	REST    *rest.RESTClient
	Dynamic dynamic.Interface
}
