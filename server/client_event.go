package server

import (
	"context"

	eventv1alpha1 "github.com/kristofferahl/aeto/apis/event/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

func (c *AetoClient) NewEventV1Alpha1Client() (*rest.RESTClient, error) {
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

func (c *AetoClient) EventV1Alpha1(namespace string) EventV1Alpha1 {
	return &eventv1Alpha1{
		restClient: c.eventv1Alpha1,
		ns:         namespace,
	}
}

type EventV1Alpha1 interface {
	ListEventStreamChunks() (*eventv1alpha1.EventStreamChunkList, error)
	GetEventStreamChunk(name string) (*eventv1alpha1.EventStreamChunk, error)
}

type eventv1Alpha1 struct {
	restClient rest.Interface
	ns         string
}

func (c *eventv1Alpha1) ListEventStreamChunks() (*eventv1alpha1.EventStreamChunkList, error) {
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

func (c *eventv1Alpha1) GetEventStreamChunk(name string) (*eventv1alpha1.EventStreamChunk, error) {
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
