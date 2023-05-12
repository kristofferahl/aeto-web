package server

import (
	"context"

	corev1alpha1 "github.com/kristofferahl/aeto/apis/core/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

type AetoClient struct {
	restConfig   *rest.Config
	coreV1Alpha1 *rest.RESTClient
}

func NewForConfig(c *rest.Config) (*AetoClient, error) {
	client := &AetoClient{
		restConfig: c,
	}

	corev1Alpha1Client, _ := client.NewCoreV1Alpha1Client()
	return &AetoClient{
		restConfig:   c,
		coreV1Alpha1: corev1Alpha1Client,
	}, nil
}

func (c *AetoClient) NewCoreV1Alpha1Client() (*rest.RESTClient, error) {
	config := *c.restConfig
	config.ContentConfig.GroupVersion = &corev1alpha1.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *AetoClient) CoreV1Alpha1(namespace string) ProjectInterface {
	return &corev1Alpha1Client{
		restClient: c.coreV1Alpha1,
		ns:         namespace,
	}
}

type ProjectInterface interface {
	List(opts metav1.ListOptions) (*corev1alpha1.TenantList, error)
}

type corev1Alpha1Client struct {
	restClient rest.Interface
	ns         string
}

func (c *corev1Alpha1Client) List(opts metav1.ListOptions) (*corev1alpha1.TenantList, error) {
	result := corev1alpha1.TenantList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("tenants").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}
