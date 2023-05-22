package server

import (
	"context"

	corev1alpha1 "github.com/kristofferahl/aeto/apis/core/v1alpha1"
	eventv1alpha1 "github.com/kristofferahl/aeto/apis/event/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
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

	corev1Alpha1Client, err := client.newCoreV1Alpha1Client()
	if err != nil {
		return nil, err
	}

	eventv1Alpha1Client, err := client.newEventV1Alpha1Client()
	if err != nil {
		return nil, err
	}

	return &AetoClient{
		restConfig:    c,
		corev1Alpha1:  corev1Alpha1Client,
		eventv1Alpha1: eventv1Alpha1Client,
	}, nil
}

func (c *AetoClient) newCoreV1Alpha1Client() (*rest.RESTClient, error) {
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

func (c *AetoClient) newEventV1Alpha1Client() (*rest.RESTClient, error) {
	config := *c.restConfig
	config.ContentConfig.GroupVersion = &eventv1alpha1.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *AetoClient) CoreV1Alpha1(namespace string) CoreV1Alpha1Interface {
	return &corev1Alpha1Client{
		restClient: c.corev1Alpha1,
		ns:         namespace,
	}
}

func (c *AetoClient) EventV1Alpha1(namespace string) EventV1Alpha1Interface {
	return &eventv1Alpha1Client{
		restClient: c.eventv1Alpha1,
		ns:         namespace,
	}
}

type CoreV1Alpha1Interface interface {
	ListTenants(opts metav1.ListOptions) (*corev1alpha1.TenantList, error)
	GetTenant(name string) (*corev1alpha1.Tenant, error)
	ListBlueprints(opts metav1.ListOptions) (*corev1alpha1.BlueprintList, error)
	GetBlueprint(name string) (*corev1alpha1.Blueprint, error)
	ListResourceSets(opts metav1.ListOptions) (*corev1alpha1.ResourceSetList, error)
	GetResourceSet(name string) (*corev1alpha1.ResourceSet, error)
	ListResourceTemplates(opts metav1.ListOptions) (*corev1alpha1.ResourceTemplateList, error)
	GetResourceTemplate(name string) (*corev1alpha1.ResourceTemplate, error)
}

type EventV1Alpha1Interface interface {
	ListEventStreamChunks(opts metav1.ListOptions) (*eventv1alpha1.EventStreamChunkList, error)
	GetEventStreamChunk(name string) (*eventv1alpha1.EventStreamChunk, error)
}

type corev1Alpha1Client struct {
	restClient rest.Interface
	ns         string
}

type eventv1Alpha1Client struct {
	restClient rest.Interface
	ns         string
}

func (c *corev1Alpha1Client) ListTenants(opts metav1.ListOptions) (*corev1alpha1.TenantList, error) {
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

func (c *corev1Alpha1Client) GetTenant(name string) (*corev1alpha1.Tenant, error) {
	result := corev1alpha1.Tenant{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Name(name).
		Resource("tenants").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}
func (c *corev1Alpha1Client) ListBlueprints(opts metav1.ListOptions) (*corev1alpha1.BlueprintList, error) {
	result := corev1alpha1.BlueprintList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("blueprints").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *corev1Alpha1Client) GetBlueprint(name string) (*corev1alpha1.Blueprint, error) {
	result := corev1alpha1.Blueprint{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Name(name).
		Resource("blueprints").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *corev1Alpha1Client) ListResourceSets(opts metav1.ListOptions) (*corev1alpha1.ResourceSetList, error) {
	result := corev1alpha1.ResourceSetList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("resourcesets").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *corev1Alpha1Client) GetResourceSet(name string) (*corev1alpha1.ResourceSet, error) {
	result := corev1alpha1.ResourceSet{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Name(name).
		Resource("resourcesets").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *corev1Alpha1Client) ListResourceTemplates(opts metav1.ListOptions) (*corev1alpha1.ResourceTemplateList, error) {
	result := corev1alpha1.ResourceTemplateList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("resourcetemplates").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *corev1Alpha1Client) GetResourceTemplate(name string) (*corev1alpha1.ResourceTemplate, error) {
	result := corev1alpha1.ResourceTemplate{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Name(name).
		Resource("resourcetemplates").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *eventv1Alpha1Client) ListEventStreamChunks(opts metav1.ListOptions) (*eventv1alpha1.EventStreamChunkList, error) {
	result := eventv1alpha1.EventStreamChunkList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("eventstreamchunks").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *eventv1Alpha1Client) GetEventStreamChunk(name string) (*eventv1alpha1.EventStreamChunk, error) {
	result := eventv1alpha1.EventStreamChunk{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Name(name).
		Resource("eventstreamchunks").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}
