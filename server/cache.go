package server

import (
	"fmt"

	corev1alpha1 "github.com/kristofferahl/aeto/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	cache = &LocalCache{
		tenant: &Cache[corev1alpha1.Tenant]{
			data: make(map[string]corev1alpha1.Tenant),
		},
		blueprint: &Cache[corev1alpha1.Blueprint]{
			data: make(map[string]corev1alpha1.Blueprint),
		},
		resourceSets: &Cache[corev1alpha1.ResourceSet]{
			data: make(map[string]corev1alpha1.ResourceSet),
		},
		resourceTemplates: &Cache[corev1alpha1.ResourceTemplate]{
			data: make(map[string]corev1alpha1.ResourceTemplate),
		},
	}
)

type LocalCache struct {
	tenant            ResourceCache[corev1alpha1.Tenant]
	blueprint         ResourceCache[corev1alpha1.Blueprint]
	resourceSets      ResourceCache[corev1alpha1.ResourceSet]
	resourceTemplates ResourceCache[corev1alpha1.ResourceTemplate]
}

type ResourceCache[T any] interface {
	Add(types.UID, T)
	Update(types.UID, T)
	Delete(types.UID)
	Items(filters ...func(i T) bool) []T
}

type Cache[T any] struct {
	data map[string]T
}

func (s Cache[T]) Add(id types.UID, obj T) {
	if id == "" {
		panic(fmt.Errorf("id must not be empty"))
	}
	s.data[string(id)] = obj
}

func (s Cache[T]) Update(id types.UID, obj T) {
	if id == "" {
		panic(fmt.Errorf("id must not be empty"))
	}
	s.data[string(id)] = obj
}

func (s Cache[T]) Delete(id types.UID) {
	if id == "" {
		panic(fmt.Errorf("id must not be empty"))
	}
	delete(s.data, string(id))
}

func (s Cache[T]) Items(filters ...func(i T) bool) []T {
	r := make([]T, 0)
	for _, v := range s.data {
		match := true
		for _, f := range filters {
			if !f(v) {
				match = false
				break
			}
		}
		if match {
			r = append(r, v)
		}
	}
	return r
}
