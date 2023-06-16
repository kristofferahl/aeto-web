package server

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	corev1alpha1 "github.com/kristofferahl/aeto/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	k8scache "k8s.io/client-go/tools/cache"
)

type CoreV1Alpha1Client struct {
	REST    *rest.RESTClient
	Dynamic dynamic.Interface
}

func (c *AetoClient) NewCoreV1Alpha1Client() (*CoreV1Alpha1Client, error) {
	config := *c.restConfig
	config.ContentConfig.GroupVersion = &corev1alpha1.GroupVersion
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

	return &CoreV1Alpha1Client{
		REST:    client,
		Dynamic: dynamicClient,
	}, nil
}

func (c *AetoClient) CoreV1Alpha1(namespace string) CoreV1Alpha1 {
	return &corev1Alpha1{
		client: c.corev1Alpha1,
		ns:     namespace,
	}
}

type CoreV1Alpha1 interface {
	Watch() error
	ListTenants(filters ...func(i corev1alpha1.Tenant) bool) (*corev1alpha1.TenantList, error)
	GetTenant(name string) (*corev1alpha1.Tenant, error)
	ListBlueprints(filters ...func(i corev1alpha1.Blueprint) bool) (*corev1alpha1.BlueprintList, error)
	GetBlueprint(name string) (*corev1alpha1.Blueprint, error)
	ListResourceSets(filters ...func(i corev1alpha1.ResourceSet) bool) (*corev1alpha1.ResourceSetList, error)
	GetResourceSet(name string) (*corev1alpha1.ResourceSet, error)
	ListResourceTemplates(filters ...func(i corev1alpha1.ResourceTemplate) bool) (*corev1alpha1.ResourceTemplateList, error)
	GetResourceTemplate(name string) (*corev1alpha1.ResourceTemplate, error)
}

type corev1Alpha1 struct {
	client *CoreV1Alpha1Client
	ns     string
}

func (c *corev1Alpha1) Watch() error {
	if err := Watch(
		corev1alpha1.GroupVersion.WithResource("tenants"),
		c.client.Dynamic,
		func() corev1alpha1.Tenant {
			return corev1alpha1.Tenant{}
		},
		cache.tenant); err != nil {
		return err
	}

	if err := Watch(
		corev1alpha1.GroupVersion.WithResource("blueprints"),
		c.client.Dynamic,
		func() corev1alpha1.Blueprint {
			return corev1alpha1.Blueprint{}
		},
		cache.blueprint); err != nil {
		return err
	}

	if err := Watch(
		corev1alpha1.GroupVersion.WithResource("resourcesets"),
		c.client.Dynamic,
		func() corev1alpha1.ResourceSet {
			return corev1alpha1.ResourceSet{}
		},
		cache.resourceSets); err != nil {
		return err
	}

	if err := NewWatcher(corev1alpha1.GroupVersion.WithResource("resourcetemplates"), c.client.Dynamic, k8scache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			log.Println("Add", "resourcetemplates", u.GetUID())
			result := corev1alpha1.ResourceTemplate{}
			err := c.client.REST.
				Get().
				Namespace(c.ns).
				Name(u.GetName()).
				Resource("resourcetemplates").
				Do(context.Background()).
				Into(&result)
			if err != nil {
				log.Println("Add", "resourcetemplates", fmt.Sprintf("error fetching resource %s/%s, err:", c.ns, u.GetName()), err)
				return
			}
			cache.resourceTemplates.Add(u.GetUID(), result)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			ou := oldObj.(*unstructured.Unstructured)
			nu := newObj.(*unstructured.Unstructured)
			log.Println("Update", "resourcetemplates", ou.GetResourceVersion(), nu.GetResourceVersion())
			result := corev1alpha1.ResourceTemplate{}
			err := c.client.REST.
				Get().
				Namespace(c.ns).
				Name(nu.GetName()).
				Resource("resourcetemplates").
				Do(context.Background()).
				Into(&result)
			if err != nil {
				log.Println("Add", "resourcetemplates", fmt.Sprintf("error fetching resource %s/%s, err:", c.ns, nu.GetName()), err)
				return
			}
			cache.resourceTemplates.Update(nu.GetUID(), result)
		},
		DeleteFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			log.Println("Delete", "resourcetemplates", u.GetUID())
			cache.resourceTemplates.Delete(u.GetUID())
		},
	}); err != nil {
		return err
	}

	return nil
}

func (c *corev1Alpha1) ListTenants(filters ...func(i corev1alpha1.Tenant) bool) (*corev1alpha1.TenantList, error) {
	result := corev1alpha1.TenantList{}

	filters = append(filters, func(i corev1alpha1.Tenant) bool {
		return i.GetNamespace() == c.ns
	})
	result.Items = cache.tenant.Items(filters...)

	sort.Slice(result.Items, func(i, j int) bool {
		return strings.Compare(result.Items[i].NamespacedName().String(), result.Items[j].NamespacedName().String()) == -1
	})

	return &result, nil
}

func (c *corev1Alpha1) GetTenant(name string) (*corev1alpha1.Tenant, error) {
	result, err := c.ListTenants(func(i corev1alpha1.Tenant) bool {
		return i.Name == name
	})
	return one(result.Items, err, &corev1alpha1.Tenant{})
}

func (c *corev1Alpha1) ListBlueprints(filters ...func(i corev1alpha1.Blueprint) bool) (*corev1alpha1.BlueprintList, error) {
	result := corev1alpha1.BlueprintList{}

	filters = append(filters, func(i corev1alpha1.Blueprint) bool {
		return i.Namespace == c.ns
	})
	result.Items = cache.blueprint.Items(filters...)

	sort.Slice(result.Items, func(i, j int) bool {
		return strings.Compare(result.Items[i].NamespacedName().String(), result.Items[j].NamespacedName().String()) == -1
	})

	return &result, nil
}

func (c *corev1Alpha1) GetBlueprint(name string) (*corev1alpha1.Blueprint, error) {
	result, err := c.ListBlueprints(func(i corev1alpha1.Blueprint) bool {
		return i.Name == name
	})
	return one(result.Items, err, &corev1alpha1.Blueprint{})
}

func (c *corev1Alpha1) ListResourceSets(filters ...func(i corev1alpha1.ResourceSet) bool) (*corev1alpha1.ResourceSetList, error) {
	result := corev1alpha1.ResourceSetList{}

	filters = append(filters, func(i corev1alpha1.ResourceSet) bool {
		return i.GetNamespace() == c.ns
	})
	result.Items = cache.resourceSets.Items(filters...)

	sort.Slice(result.Items, func(i, j int) bool {
		return strings.Compare(result.Items[i].NamespacedName().String(), result.Items[j].NamespacedName().String()) == -1
	})

	return &result, nil
}

func (c *corev1Alpha1) GetResourceSet(name string) (*corev1alpha1.ResourceSet, error) {
	result, err := c.ListResourceSets(func(i corev1alpha1.ResourceSet) bool {
		return i.Name == name
	})
	return one(result.Items, err, &corev1alpha1.ResourceSet{})
}

func (c *corev1Alpha1) ListResourceTemplates(filters ...func(i corev1alpha1.ResourceTemplate) bool) (*corev1alpha1.ResourceTemplateList, error) {
	result := corev1alpha1.ResourceTemplateList{}

	filters = append(filters, func(i corev1alpha1.ResourceTemplate) bool {
		return i.GetNamespace() == c.ns
	})
	result.Items = cache.resourceTemplates.Items(filters...)

	sort.Slice(result.Items, func(i, j int) bool {
		return strings.Compare(result.Items[i].NamespacedName().String(), result.Items[j].NamespacedName().String()) == -1
	})

	return &result, nil
}

func (c *corev1Alpha1) GetResourceTemplate(name string) (*corev1alpha1.ResourceTemplate, error) {
	result, err := c.ListResourceTemplates(func(i corev1alpha1.ResourceTemplate) bool {
		return i.Name == name
	})
	return one(result.Items, err, &corev1alpha1.ResourceTemplate{})
}
