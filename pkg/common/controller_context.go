package common

import (
	"math/rand"
	"time"

	"k8s.io/client-go/informers"
)

const (
	minResyncPeriod = 20 * time.Minute
)

func resyncPeriod() func() time.Duration {
	return func() time.Duration {
		factor := rand.Float64() + 1
		return time.Duration(float64(minResyncPeriod.Nanoseconds()) * factor)
	}
}

// ControllerContext stores all the informers for a variety of kubernetes objects.
type ControllerContext struct {
	ClientBuilder                 *Builder
	KubeNamespacedInformerFactory informers.SharedInformerFactory

	Stop <-chan struct{}

	InformersStarted chan struct{}

	ResyncPeriod func() time.Duration
}

// CreateControllerContext creates the ControllerContext with the ClientBuilder.
func CreateControllerContext(cb *Builder, stop <-chan struct{}, targetNamespace string) *ControllerContext {
	kubeClient := cb.KubeClientOrDie("kube-shared-informer")
	kubeNamespacedSharedInformer := informers.NewFilteredSharedInformerFactory(kubeClient, resyncPeriod()(), targetNamespace, nil)

	return &ControllerContext{
		ClientBuilder:                 cb,
		KubeNamespacedInformerFactory: kubeNamespacedSharedInformer,
		Stop:                          stop,
		InformersStarted:              make(chan struct{}),
		ResyncPeriod:                  resyncPeriod(),
	}
}
