package server

import (
	"encoding/json"
	"fmt"
	"time"

	acmawsv1alpa1 "github.com/kristofferahl/aeto/apis/acm.aws/v1alpha1"
	corev1alpha1 "github.com/kristofferahl/aeto/apis/core/v1alpha1"
	eventv1alpha1 "github.com/kristofferahl/aeto/apis/event/v1alpha1"
	route53awsv1alpha1 "github.com/kristofferahl/aeto/apis/route53.aws/v1alpha1"
	sustainabilityv1alpha1 "github.com/kristofferahl/aeto/apis/sustainability/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	cache = &InMemoryCache{
		acmawsCertificate: &Cache[acmawsv1alpa1.Certificate]{
			data: make(map[string]CacheEntry[acmawsv1alpa1.Certificate]),
		},
		acmawsCertificateConnector: &Cache[acmawsv1alpa1.CertificateConnector]{
			data: make(map[string]CacheEntry[acmawsv1alpa1.CertificateConnector]),
		},

		coreTenant: &Cache[corev1alpha1.Tenant]{
			data: make(map[string]CacheEntry[corev1alpha1.Tenant]),
		},
		coreBlueprint: &Cache[corev1alpha1.Blueprint]{
			data: make(map[string]CacheEntry[corev1alpha1.Blueprint]),
		},
		coreResourceSet: &Cache[corev1alpha1.ResourceSet]{
			data: make(map[string]CacheEntry[corev1alpha1.ResourceSet]),
		},
		coreResourceTemplate: &Cache[corev1alpha1.ResourceTemplate]{
			data: make(map[string]CacheEntry[corev1alpha1.ResourceTemplate]),
		},

		eventEventStreamChunk: &Cache[eventv1alpha1.EventStreamChunk]{
			data: make(map[string]CacheEntry[eventv1alpha1.EventStreamChunk]),
		},

		route53awsHostedZone: &Cache[route53awsv1alpha1.HostedZone]{
			data: make(map[string]CacheEntry[route53awsv1alpha1.HostedZone]),
		},

		sustainabilitySavingsPolicy: &Cache[sustainabilityv1alpha1.SavingsPolicy]{
			data: make(map[string]CacheEntry[sustainabilityv1alpha1.SavingsPolicy]),
		},
	}
)

type InMemoryCache struct {
	acmawsCertificate          ResourceCache[acmawsv1alpa1.Certificate]
	acmawsCertificateConnector ResourceCache[acmawsv1alpa1.CertificateConnector]

	coreTenant           ResourceCache[corev1alpha1.Tenant]
	coreBlueprint        ResourceCache[corev1alpha1.Blueprint]
	coreResourceSet      ResourceCache[corev1alpha1.ResourceSet]
	coreResourceTemplate ResourceCache[corev1alpha1.ResourceTemplate]

	eventEventStreamChunk ResourceCache[eventv1alpha1.EventStreamChunk]

	route53awsHostedZone ResourceCache[route53awsv1alpha1.HostedZone]

	sustainabilitySavingsPolicy ResourceCache[sustainabilityv1alpha1.SavingsPolicy]
}

type CacheEvent struct {
	Timestamp string      `json:"ts"`
	EventType string      `json:"type"`
	Resource  interface{} `json:"resource"`
}

func (e *CacheEvent) Payload() ([]byte, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	fb, err := ApplyResourceFilter("ts,type,resource(%s)", b)
	if err != nil {
		return nil, err
	}
	return fb, err
}

type ResourceCache[T CacheableEntry] interface {
	Add(id types.UID, version string, obj T)
	Update(id types.UID, newVersion string, obj T)
	Delete(id types.UID)
	Items(filters ...func(i T) bool) []T
}

type CacheableEntry interface {
	NamespacedName() types.NamespacedName
}

type Cache[T CacheableEntry] struct {
	data map[string]CacheEntry[T]
}

type CacheEntry[T CacheableEntry] struct {
	Version  string
	Resource T
}

func (s Cache[T]) Add(id types.UID, version string, obj T) {
	if id == "" {
		panic(fmt.Errorf("id must not be empty"))
	}
	s.data[string(id)] = CacheEntry[T]{
		Version:  version,
		Resource: obj,
	}
	eventManager.Publish("change", &CacheEvent{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		EventType: "ResourceAdded",
		Resource:  obj,
	})
}

func (s Cache[T]) Update(id types.UID, newVersion string, obj T) {
	if id == "" {
		panic(fmt.Errorf("id must not be empty"))
	}
	oldObj := s.data[string(id)]
	if oldObj.Version != newVersion {
		s.data[string(id)] = CacheEntry[T]{
			Version:  newVersion,
			Resource: obj,
		}
		eventManager.Publish("change", &CacheEvent{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			EventType: "ResourceUpdated",
			Resource:  obj,
		})
	}
}

func (s Cache[T]) Delete(id types.UID) {
	if id == "" {
		panic(fmt.Errorf("id must not be empty"))
	}
	obj := s.data[string(id)]
	delete(s.data, string(id))
	eventManager.Publish("change", &CacheEvent{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		EventType: "ResourceDeleted",
		Resource:  obj,
	})
}

func (s Cache[T]) Items(filters ...func(i T) bool) []T {
	r := make([]T, 0)
	for _, v := range s.data {
		match := true
		for _, f := range filters {
			if !f(v.Resource) {
				match = false
				break
			}
		}
		if match {
			r = append(r, v.Resource)
		}
	}
	return r
}
