package server

import (
	"sort"
	"strings"

	eventv1alpha1 "github.com/kristofferahl/aeto/apis/event/v1alpha1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

func (c *AetoClient) NewEventV1Alpha1Client() (*KubernetesClient, error) {
	config := *c.restConfig
	config.ContentConfig.GroupVersion = &eventv1alpha1.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfigAndClient(&config, client.Client)
	if err != nil {
		return nil, err
	}

	return &KubernetesClient{
		REST:    client,
		Dynamic: dynamicClient,
	}, nil
}

func (c *AetoClient) EventV1Alpha1(namespace string) EventV1Alpha1 {
	return &eventv1Alpha1{
		client: c.eventv1Alpha1,
		ns:     namespace,
	}
}

type EventV1Alpha1 interface {
	Watch() error
	ListEventStreamChunks(filters ...func(i eventv1alpha1.EventStreamChunk) bool) (*eventv1alpha1.EventStreamChunkList, error)
	GetEventStreamChunk(name string) (*eventv1alpha1.EventStreamChunk, error)
}

type eventv1Alpha1 struct {
	client *KubernetesClient
	ns     string
}

func (c *eventv1Alpha1) Watch() error {
	if err := Watch(
		eventv1alpha1.GroupVersion.WithResource("eventstreamchunks"),
		c.client.Dynamic,
		func() eventv1alpha1.EventStreamChunk {
			return eventv1alpha1.EventStreamChunk{}
		},
		cache.eventEventStreamChunk); err != nil {
		return err
	}
	return nil
}

func (c *eventv1Alpha1) ListEventStreamChunks(filters ...func(i eventv1alpha1.EventStreamChunk) bool) (*eventv1alpha1.EventStreamChunkList, error) {
	result := eventv1alpha1.EventStreamChunkList{}

	filters = append(filters, func(i eventv1alpha1.EventStreamChunk) bool {
		return i.GetNamespace() == c.ns
	})
	result.Items = cache.eventEventStreamChunk.Items(filters...)

	sort.Slice(result.Items, func(i, j int) bool {
		return strings.Compare(result.Items[i].NamespacedName().String(), result.Items[j].NamespacedName().String()) == -1
	})

	return &result, nil
}

func (c *eventv1Alpha1) GetEventStreamChunk(name string) (*eventv1alpha1.EventStreamChunk, error) {
	result, err := c.ListEventStreamChunks(func(i eventv1alpha1.EventStreamChunk) bool {
		return i.Name == name
	})
	return one(result.Items, err, &eventv1alpha1.EventStreamChunk{})
}
