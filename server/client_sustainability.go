package server

import (
	"sort"
	"strings"

	sustainabilityv1alpha1 "github.com/kristofferahl/aeto/apis/sustainability/v1alpha1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

func (c *AetoClient) NewSustainabilityV1Alpha1Client() (*KubernetesClient, error) {
	config := *c.restConfig
	config.ContentConfig.GroupVersion = &sustainabilityv1alpha1.GroupVersion
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

func (c *AetoClient) SustainabilityV1Alpha1(namespace string) SustainabilityV1Alpha1 {
	return &sustainabilityV1Alpha1{
		client: c.sustainabilityv1Alpha1,
		ns:     namespace,
	}
}

type SustainabilityV1Alpha1 interface {
	Watch() error
	ListSavingsPolicies(filters ...func(i sustainabilityv1alpha1.SavingsPolicy) bool) (*sustainabilityv1alpha1.SavingsPolicyList, error)
	GetSavingsPolicy(name string) (*sustainabilityv1alpha1.SavingsPolicy, error)
}

type sustainabilityV1Alpha1 struct {
	client *KubernetesClient
	ns     string
}

func (c *sustainabilityV1Alpha1) Watch() error {
	if err := Watch(
		sustainabilityv1alpha1.GroupVersion.WithResource("savingspolicies"),
		c.client.Dynamic,
		func() sustainabilityv1alpha1.SavingsPolicy {
			return sustainabilityv1alpha1.SavingsPolicy{}
		},
		cache.sustainabilitySavingsPolicy); err != nil {
		return err
	}
	return nil
}

func (c *sustainabilityV1Alpha1) ListSavingsPolicies(filters ...func(i sustainabilityv1alpha1.SavingsPolicy) bool) (*sustainabilityv1alpha1.SavingsPolicyList, error) {
	result := sustainabilityv1alpha1.SavingsPolicyList{}

	filters = append(filters, func(i sustainabilityv1alpha1.SavingsPolicy) bool {
		return i.GetNamespace() == c.ns
	})
	result.Items = cache.sustainabilitySavingsPolicy.Items(filters...)

	sort.Slice(result.Items, func(i, j int) bool {
		return strings.Compare(result.Items[i].NamespacedName().String(), result.Items[j].NamespacedName().String()) == -1
	})

	return &result, nil
}

func (c *sustainabilityV1Alpha1) GetSavingsPolicy(name string) (*sustainabilityv1alpha1.SavingsPolicy, error) {
	result, err := c.ListSavingsPolicies(func(i sustainabilityv1alpha1.SavingsPolicy) bool {
		return i.Name == name
	})
	return one(result.Items, err, &sustainabilityv1alpha1.SavingsPolicy{})
}
