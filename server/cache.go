package server

import (
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/kristofferahl/aeto-web/server/sse"
	corev1alpha1 "github.com/kristofferahl/aeto/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	cache = &InMemoryCache{
		changestream: &ChangeStream{
			recordAfter: time.Now().UTC().Add(1 * time.Minute),
			maxEvents:   10,
			events:      make([]CacheEvent, 0),
		},
		tenant: &Cache[corev1alpha1.Tenant]{
			data: make(map[string]CacheEntry[corev1alpha1.Tenant]),
		},
		blueprint: &Cache[corev1alpha1.Blueprint]{
			data: make(map[string]CacheEntry[corev1alpha1.Blueprint]),
		},
		resourceSets: &Cache[corev1alpha1.ResourceSet]{
			data: make(map[string]CacheEntry[corev1alpha1.ResourceSet]),
		},
		resourceTemplates: &Cache[corev1alpha1.ResourceTemplate]{
			data: make(map[string]CacheEntry[corev1alpha1.ResourceTemplate]),
		},
	}
)

type InMemoryCache struct {
	changestream      *ChangeStream
	tenant            ResourceCache[corev1alpha1.Tenant]
	blueprint         ResourceCache[corev1alpha1.Blueprint]
	resourceSets      ResourceCache[corev1alpha1.ResourceSet]
	resourceTemplates ResourceCache[corev1alpha1.ResourceTemplate]
}

type ChangeStream struct {
	recordAfter time.Time
	maxEvents   int
	events      []CacheEvent
}

func (s *ChangeStream) AddEvent(e CacheEvent) {
	now := time.Now().UTC()
	if now.After(s.recordAfter) {
		e.time = now
		e.Timestamp = e.time.Format(time.RFC3339)
		s.events = append(s.events, e)
		if len(s.events) > s.maxEvents {
			s.events = s.events[len(s.events)-s.maxEvents:] // Remove older events
		}
	}
}

func (s *ChangeStream) TakeLast(n int) []CacheEvent {
	changes := cache.changestream
	sort.Slice(changes.events, func(i, j int) bool {
		return changes.events[i].time.Before(changes.events[j].time)
	})
	ne := len(changes.events)
	if n > ne {
		n = ne
	}
	return changes.events[ne-n:]
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

type CacheEvent struct {
	Type      string `json:"type"`
	Resource  string `json:"resource"`
	Change    string `json:"change"`
	time      time.Time
	Timestamp string `json:"ts"`
}

func (s Cache[T]) Add(id types.UID, version string, obj T) {
	if id == "" {
		panic(fmt.Errorf("id must not be empty"))
	}
	s.data[string(id)] = CacheEntry[T]{
		Version:  version,
		Resource: obj,
	}
	cache.changestream.AddEvent(CacheEvent{
		Change:   "Added",
		Type:     reflect.TypeOf(obj).Name(),
		Resource: obj.NamespacedName().String(),
	})
	eventManager.Publish("change", sse.Event{
		Type:    "ResourceAdded",
		Payload: obj,
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
		cache.changestream.AddEvent(CacheEvent{
			Change:   "Updated",
			Type:     reflect.TypeOf(obj).Name(),
			Resource: obj.NamespacedName().String(),
		})
		eventManager.Publish("change", sse.Event{
			Type:    "ResourceChanged",
			Payload: obj,
		})
	}
}

func (s Cache[T]) Delete(id types.UID) {
	if id == "" {
		panic(fmt.Errorf("id must not be empty"))
	}
	obj := s.data[string(id)]
	delete(s.data, string(id))
	cache.changestream.AddEvent(CacheEvent{
		Change:   "Deleted",
		Type:     reflect.TypeOf(obj).Name(),
		Resource: obj.Resource.NamespacedName().String(),
	})
	eventManager.Publish("change", sse.Event{
		Type:    "ResourceDeleted",
		Payload: obj,
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
