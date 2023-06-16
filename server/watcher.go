package server

import (
	"fmt"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	k8scache "k8s.io/client-go/tools/cache"
)

func Watch[T CacheableEntry](resource schema.GroupVersionResource, client dynamic.Interface, resourceFactory func() T, resourceCache ResourceCache[T]) error {
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, time.Minute*5, "", nil)

	informer := factory.ForResource(resource).Informer()
	informer.AddEventHandler(k8scache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			log.Println("Add", resource.Resource, u.GetUID())
			r := resourceFactory()
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), &r)
			if err != nil {
				log.Println(fmt.Sprintf("error converting to type %s", resource.Resource), err)
			}
			resourceCache.Add(u.GetUID(), u.GetResourceVersion(), r)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			ou := oldObj.(*unstructured.Unstructured)
			nu := newObj.(*unstructured.Unstructured)
			log.Println("Update", resource.Resource, ou.GetResourceVersion(), nu.GetResourceVersion())
			r := resourceFactory()
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(nu.UnstructuredContent(), &r)
			if err != nil {
				log.Println(fmt.Sprintf("error converting to type %s", resource.Resource), err)
			}
			resourceCache.Update(nu.GetUID(), nu.GetResourceVersion(), r)
		},
		DeleteFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			log.Println("Delete", resource.Resource, u.GetUID())
			resourceCache.Delete(u.GetUID())
		},
	})

	log.Println("Watching", resource.Resource)
	// TODO: Stop the informer before exiting program
	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer cancel()
	go informer.Run(make(<-chan struct{}))

	return nil
}

func NewWatcher(resource schema.GroupVersionResource, client dynamic.Interface, resourceEventHandler k8scache.ResourceEventHandlerFuncs) error {
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, time.Minute*5, "", nil)

	informer := factory.ForResource(resource).Informer()
	informer.AddEventHandler(resourceEventHandler)

	log.Println("Watching", resource.Resource)
	// TODO: Stop the informer before exiting program
	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer cancel()
	go informer.Run(make(<-chan struct{}))

	return nil
}
