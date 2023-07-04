package server

import (
	"sort"
	"strings"

	route53awsv1alpha1 "github.com/kristofferahl/aeto/apis/route53.aws/v1alpha1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

func (c *AetoClient) NewRoute53AwsV1Alpha1Client() (*KubernetesClient, error) {
	config := *c.restConfig
	config.ContentConfig.GroupVersion = &route53awsv1alpha1.GroupVersion
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

func (c *AetoClient) Route53AwsV1Alpha1(namespace string) Route53AwsV1Alpha1 {
	return &route53AwsV1Alpha1{
		client: c.route53AwsV1Alpha1,
		ns:     namespace,
	}
}

type Route53AwsV1Alpha1 interface {
	Watch() error
	ListHostedZones(filters ...func(i route53awsv1alpha1.HostedZone) bool) (*route53awsv1alpha1.HostedZoneList, error)
	GetHostedZone(name string) (*route53awsv1alpha1.HostedZone, error)
}

type route53AwsV1Alpha1 struct {
	client *KubernetesClient
	ns     string
}

func (c *route53AwsV1Alpha1) Watch() error {
	if err := Watch(
		route53awsv1alpha1.GroupVersion.WithResource("hostedzones"),
		c.client.Dynamic,
		func() route53awsv1alpha1.HostedZone {
			return route53awsv1alpha1.HostedZone{}
		},
		cache.route53awsHostedZone); err != nil {
		return err
	}
	return nil
}

func (c *route53AwsV1Alpha1) ListHostedZones(filters ...func(i route53awsv1alpha1.HostedZone) bool) (*route53awsv1alpha1.HostedZoneList, error) {
	result := route53awsv1alpha1.HostedZoneList{}

	filters = append(filters, func(i route53awsv1alpha1.HostedZone) bool {
		return i.GetNamespace() == c.ns
	})
	result.Items = cache.route53awsHostedZone.Items(filters...)

	sort.Slice(result.Items, func(i, j int) bool {
		return strings.Compare(result.Items[i].NamespacedName().String(), result.Items[j].NamespacedName().String()) == -1
	})

	return &result, nil
}

func (c *route53AwsV1Alpha1) GetHostedZone(name string) (*route53awsv1alpha1.HostedZone, error) {
	result, err := c.ListHostedZones(func(i route53awsv1alpha1.HostedZone) bool {
		return i.Name == name
	})
	return one(result.Items, err, &route53awsv1alpha1.HostedZone{})
}
