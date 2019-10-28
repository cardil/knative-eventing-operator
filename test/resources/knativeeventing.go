package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/openshift-knative/knative-eventing-operator/pkg/apis/eventing/v1alpha1"
	eventingv1alpha1 "github.com/openshift-knative/knative-eventing-operator/pkg/client/clientset/versioned/typed/eventing/v1alpha1"
	"github.com/openshift-knative/knative-eventing-operator/test"
	"github.com/pkg/errors"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"knative.dev/pkg/test/logging"
)

const (
	// Interval specifies the time between two polls.
	Interval = 10 * time.Second
	// Timeout specifies the timeout for the function PollImmediate to reach a certain status.
	Timeout = 5 * time.Minute
)

// WaitForKnativeEventingState polls the status of the KnativeServing called name
// from client every `interval` until `inState` returns `true` indicating it
// is done, returns an error or timeout.
func WaitForKnativeEventingState(clients eventingv1alpha1.KnativeEventingInterface, name string,
	inState func(s *v1alpha1.KnativeEventing, err error) (bool, error)) (*v1alpha1.KnativeEventing, error) {
	span := logging.GetEmitableSpan(context.Background(), fmt.Sprintf("WaitForKnativeEventingState/%s/%s", name, "KnativeServingIsReady"))
	defer span.End()

	var lastState *v1alpha1.KnativeEventing
	waitErr := wait.PollImmediate(Interval, Timeout, func() (bool, error) {
		lastState, err := clients.Get(name, metav1.GetOptions{})
		return inState(lastState, err)
	})

	if waitErr != nil {
		return lastState, errors.Wrapf(waitErr, "Eventing %s is not in desired state, got: %+v", name, lastState)
	}
	return lastState, nil
}

// CreateKnativeEventing creates a KnativeServing with the name names.KnativeServing under the namespace names.Namespace.
func CreateKnativeEventing(clients eventingv1alpha1.KnativeEventingInterface, names test.ResourceNames) (*v1alpha1.KnativeEventing, error) {
	ks := &v1alpha1.KnativeEventing{
		ObjectMeta: metav1.ObjectMeta{
			Name:      names.KnativeEventing,
			Namespace: names.Namespace,
		},
	}
	svc, err := clients.Create(ks)
	return svc, err
}

// IsKnativeEventingReady will check the status conditions of the KnativeServing and return true if the KnativeServing is ready.
func IsKnativeEventingReady(s *v1alpha1.KnativeEventing, err error) (bool, error) {
	return s.Status.IsReady(), err
}

// IsDeploymentAvailable will check the status conditions of the deployment and return true if the deployment is available.
func IsDeploymentAvailable(d *v1.Deployment) (bool, error) {
	return getDeploymentStatus(d) == "True", nil
}

func getDeploymentStatus(d *v1.Deployment) corev1.ConditionStatus {
	for _, dc := range d.Status.Conditions {
		if dc.Type == "Available" {
			return dc.Status
		}
	}
	return "unknown"
}
