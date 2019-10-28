// This file contains an object which encapsulates k8s clients which are useful for e2e tests.

package test

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/openshift-knative/knative-eventing-operator/pkg/client/clientset/versioned"
	eventingv1alpha1 "github.com/openshift-knative/knative-eventing-operator/pkg/client/clientset/versioned/typed/eventing/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/test"
)

// Clients holds instances of interfaces for making requests to Knative Serving.
type Clients struct {
	KubeClient *test.KubeClient
	Dynamic    dynamic.Interface
	Eventing   eventingv1alpha1.EventingV1alpha1Interface
	Config     *rest.Config
}

// NewClients instantiates and returns several clientsets required for making request to the
// Knative Serving cluster specified by the combination of clusterName and configPath.
func NewClients(configPath string, clusterName string) (*Clients, error) {
	clients := &Clients{}
	cfg, err := buildClientConfig(configPath, clusterName)
	if err != nil {
		return nil, err
	}

	// We poll, so set our limits high.
	cfg.QPS = 100
	cfg.Burst = 200

	clients.KubeClient, err = test.NewKubeClient(configPath, clusterName)
	if err != nil {
		return nil, err
	}

	clients.Dynamic, err = dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	clients.Eventing, err = newKnativeEventingAlphaClients(cfg)
	if err != nil {
		return nil, err
	}

	clients.Config = cfg
	return clients, nil
}

func buildClientConfig(kubeConfigPath string, clusterName string) (*rest.Config, error) {
	overrides := clientcmd.ConfigOverrides{}
	// Override the cluster name if provided.
	if clusterName != "" {
		overrides.Context.Cluster = clusterName
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
		&overrides).ClientConfig()
}

func newKnativeEventingAlphaClients(cfg *rest.Config) (eventingv1alpha1.EventingV1alpha1Interface, error) {
	cs, err := versioned.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cs.EventingV1alpha1(), nil
}

func (c *Clients) KnativeEventing() eventingv1alpha1.KnativeEventingInterface {
	return c.Eventing.KnativeEventings(EventingOperatorNamespace)
}

func (c *Clients) KnativeEventingAll() eventingv1alpha1.KnativeEventingInterface {
	return c.Eventing.KnativeEventings(metav1.NamespaceAll)
}
