package server

import (
	"sort"
	"strings"

	acmawsv1alpha1 "github.com/kristofferahl/aeto/apis/acm.aws/v1alpha1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

func (c *AetoClient) NewAcmAwsV1Alpha1Client() (*KubernetesClient, error) {
	config := *c.restConfig
	config.ContentConfig.GroupVersion = &acmawsv1alpha1.GroupVersion
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

func (c *AetoClient) AcmAwsV1Alpha1(namespace string) AcmAwsV1Alpha1 {
	return &acmAwsV1Alpha1{
		client: c.acmAwsV1Alpha1,
		ns:     namespace,
	}
}

type AcmAwsV1Alpha1 interface {
	Watch() error
	ListCertificates(filters ...func(i acmawsv1alpha1.Certificate) bool) (*acmawsv1alpha1.CertificateList, error)
	GetCertificate(name string) (*acmawsv1alpha1.Certificate, error)
	ListCertificateConnectors(filters ...func(i acmawsv1alpha1.CertificateConnector) bool) (*acmawsv1alpha1.CertificateConnectorList, error)
	GetCertificateConnector(name string) (*acmawsv1alpha1.CertificateConnector, error)
}

type acmAwsV1Alpha1 struct {
	client *KubernetesClient
	ns     string
}

func (c *acmAwsV1Alpha1) Watch() error {
	if err := Watch(
		acmawsv1alpha1.GroupVersion.WithResource("certificates"),
		c.client.Dynamic,
		func() acmawsv1alpha1.Certificate {
			return acmawsv1alpha1.Certificate{}
		},
		cache.acmawsCertificate); err != nil {
		return err
	}

	if err := Watch(
		acmawsv1alpha1.GroupVersion.WithResource("certificateconnectors"),
		c.client.Dynamic,
		func() acmawsv1alpha1.CertificateConnector {
			return acmawsv1alpha1.CertificateConnector{}
		},
		cache.acmawsCertificateConnector); err != nil {
		return err
	}

	return nil
}

func (c *acmAwsV1Alpha1) ListCertificates(filters ...func(i acmawsv1alpha1.Certificate) bool) (*acmawsv1alpha1.CertificateList, error) {
	result := acmawsv1alpha1.CertificateList{}

	filters = append(filters, func(i acmawsv1alpha1.Certificate) bool {
		return i.GetNamespace() == c.ns
	})
	result.Items = cache.acmawsCertificate.Items(filters...)

	sort.Slice(result.Items, func(i, j int) bool {
		return strings.Compare(result.Items[i].NamespacedName().String(), result.Items[j].NamespacedName().String()) == -1
	})

	return &result, nil
}

func (c *acmAwsV1Alpha1) GetCertificate(name string) (*acmawsv1alpha1.Certificate, error) {
	result, err := c.ListCertificates(func(i acmawsv1alpha1.Certificate) bool {
		return i.Name == name
	})
	return one(result.Items, err, &acmawsv1alpha1.Certificate{})
}

func (c *acmAwsV1Alpha1) ListCertificateConnectors(filters ...func(i acmawsv1alpha1.CertificateConnector) bool) (*acmawsv1alpha1.CertificateConnectorList, error) {
	result := acmawsv1alpha1.CertificateConnectorList{}

	filters = append(filters, func(i acmawsv1alpha1.CertificateConnector) bool {
		return i.GetNamespace() == c.ns
	})
	result.Items = cache.acmawsCertificateConnector.Items(filters...)

	sort.Slice(result.Items, func(i, j int) bool {
		return strings.Compare(result.Items[i].NamespacedName().String(), result.Items[j].NamespacedName().String()) == -1
	})

	return &result, nil
}

func (c *acmAwsV1Alpha1) GetCertificateConnector(name string) (*acmawsv1alpha1.CertificateConnector, error) {
	result, err := c.ListCertificateConnectors(func(i acmawsv1alpha1.CertificateConnector) bool {
		return i.Name == name
	})
	return one(result.Items, err, &acmawsv1alpha1.CertificateConnector{})
}
