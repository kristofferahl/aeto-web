package server

import (
	"context"

	sustainabilityv1alpha1 "github.com/kristofferahl/aeto/apis/sustainability/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

func (c *AetoClient) NewSustainabilityV1Alpha1Client() (*rest.RESTClient, error) {
	config := *c.restConfig
	config.ContentConfig.GroupVersion = &sustainabilityv1alpha1.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *AetoClient) SustainabilityV1Alpha1(namespace string) SustainabilityV1Alpha1 {
	return &sustainabilityV1Alpha1{
		restClient: c.sustainabilityv1Alpha1,
		ns:         namespace,
	}
}

type SustainabilityV1Alpha1 interface {
	ListSavingsPolicies() (*sustainabilityv1alpha1.SavingsPolicyList, error)
	GetSavingsPolicy(name string) (*sustainabilityv1alpha1.SavingsPolicy, error)
}

type sustainabilityV1Alpha1 struct {
	restClient rest.Interface
	ns         string
}

func (c *sustainabilityV1Alpha1) ListSavingsPolicies() (*sustainabilityv1alpha1.SavingsPolicyList, error) {
	result := sustainabilityv1alpha1.SavingsPolicyList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("savingspolicies").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *sustainabilityV1Alpha1) GetSavingsPolicy(name string) (*sustainabilityv1alpha1.SavingsPolicy, error) {
	result := sustainabilityv1alpha1.SavingsPolicy{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Name(name).
		Resource("savingspolicies").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}
