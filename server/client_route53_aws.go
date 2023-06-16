package server

import (
	"context"

	route53awsv1alpha1 "github.com/kristofferahl/aeto/apis/route53.aws/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

func (c *AetoClient) NewRoute53AwsV1Alpha1Client() (*rest.RESTClient, error) {
	config := *c.restConfig
	config.ContentConfig.GroupVersion = &route53awsv1alpha1.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *AetoClient) Route53AwsV1Alpha1(namespace string) Route53AwsV1Alpha1 {
	return &route53AwsV1Alpha1{
		restClient: c.route53AwsV1Alpha1,
		ns:         namespace,
	}
}

type Route53AwsV1Alpha1 interface {
	ListHostedZones() (*route53awsv1alpha1.HostedZoneList, error)
	GetHostedZone(name string) (*route53awsv1alpha1.HostedZone, error)
}

type route53AwsV1Alpha1 struct {
	restClient rest.Interface
	ns         string
}

func (c *route53AwsV1Alpha1) ListHostedZones() (*route53awsv1alpha1.HostedZoneList, error) {
	result := route53awsv1alpha1.HostedZoneList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("hostedzones").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *route53AwsV1Alpha1) GetHostedZone(name string) (*route53awsv1alpha1.HostedZone, error) {
	result := route53awsv1alpha1.HostedZone{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Name(name).
		Resource("hostedzones").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}
