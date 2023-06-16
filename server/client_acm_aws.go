package server

import (
	"context"

	acmawsv1alpha1 "github.com/kristofferahl/aeto/apis/acm.aws/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

func (c *AetoClient) NewAcmAwsV1Alpha1Client() (*rest.RESTClient, error) {
	config := *c.restConfig
	config.ContentConfig.GroupVersion = &acmawsv1alpha1.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *AetoClient) AcmAwsV1Alpha1(namespace string) AcmAwsV1Alpha1 {
	return &acmAwsV1Alpha1{
		restClient: c.acmAwsV1Alpha1,
		ns:         namespace,
	}
}

type AcmAwsV1Alpha1 interface {
	ListCertificates() (*acmawsv1alpha1.CertificateList, error)
	GetCertificate(name string) (*acmawsv1alpha1.Certificate, error)
	ListCertificateConnectors() (*acmawsv1alpha1.CertificateConnectorList, error)
	GetCertificateConnector(name string) (*acmawsv1alpha1.CertificateConnector, error)
}

type acmAwsV1Alpha1 struct {
	restClient rest.Interface
	ns         string
}

func (c *acmAwsV1Alpha1) ListCertificates() (*acmawsv1alpha1.CertificateList, error) {
	result := acmawsv1alpha1.CertificateList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("certificates").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *acmAwsV1Alpha1) GetCertificate(name string) (*acmawsv1alpha1.Certificate, error) {
	result := acmawsv1alpha1.Certificate{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Name(name).
		Resource("certificates").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *acmAwsV1Alpha1) ListCertificateConnectors() (*acmawsv1alpha1.CertificateConnectorList, error) {
	result := acmawsv1alpha1.CertificateConnectorList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("certificateconnectors").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *acmAwsV1Alpha1) GetCertificateConnector(name string) (*acmawsv1alpha1.CertificateConnector, error) {
	result := acmawsv1alpha1.CertificateConnector{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Name(name).
		Resource("certificateconnectors").
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}
