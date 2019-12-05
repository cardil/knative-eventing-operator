//+build e2e

package e2e

import (
	"github.com/cardil/operator-e2e-harness/pkg/config"
	harness "github.com/cardil/operator-e2e-harness/pkg/test"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"

	eventingv1alpha1 "github.com/openshift-knative/knative-eventing-operator/pkg/apis/eventing/v1alpha1"
	"github.com/openshift-knative/knative-eventing-operator/pkg/controller/knativeeventing"
	"github.com/openshift-knative/knative-eventing-operator/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var eventingCR = eventingv1alpha1.SchemeGroupVersion.
	WithResource("knativeeventings")

func TestBasicLifecycle(t *testing.T) {
	ctx := test.NewContext(t)
	defer ctx.Cleanup()

	t.Run("deploy eventing operator", func(t *testing.T) {
		packageName := "knative-eventing-operator"
		subscriptionName := "eventing-operator-subscription"
		_, err := harness.WithOperatorReady(ctx.Base, subscriptionName, packageName)
		if err != nil {
			t.Fatal("Failed", err)
		}
	})

	t.Run("wait until eventing cr reports ready", func(t *testing.T) {
		// FIXME: SRVKE-252 Eventing CR should'n be automatically created by operator

		var cr *unstructured.Unstructured
		var err error

		waitErr := wait.PollImmediate(config.Polling.Interval, config.Polling.Timeout, func() (bool, error) {
			ns := knativeeventing.Operand
			cr, err = ctx.Base.Clients.Dynamic.Resource(eventingCR).Namespace(ns).
				Get(ns, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			content := cr.UnstructuredContent()
			status := content["status"].(map[string]interface{})
			conditions := status["conditions"].([]interface{})
			result := false
			for _, elem := range conditions {
				condition := elem.(map[string]interface{})
				if condition["type"] == "Ready" {
					result = condition["status"] == "True"
				}
			}
			return result, nil
		})

		if waitErr != nil {
			t.Error(errors.Wrapf(
				waitErr, "waiting for eventing cr to report ready failed, got: %+v",
				cr))
		}
	})

	t.Run("check eventing pods are up", func(t *testing.T) {
		ns := knativeeventing.Operand
		pods, err := ctx.Base.Clients.KubeClient.Kube.CoreV1().
			Pods(ns).
			List(metav1.ListOptions{})
		if err != nil {
			t.Fatal(err)
		}

		size := len(pods.Items)
		if size == 0 {
			t.Errorf("In namespace %s there should be some " +
				"pods after successful subscription, but wasn't", ns)
		}
	})

	t.Run("delete eventing cr", func(t *testing.T) {
		ns := knativeeventing.Operand
		err := ctx.Base.Clients.Dynamic.Resource(eventingCR).Namespace(ns).
			Delete(ns, &metav1.DeleteOptions{})

		if err != nil {
			t.Error(err)
		}

		inState := func(pods *corev1.PodList, err error) (bool, error) {
			return len(pods.Items) == 0, err
		}

		var pods *corev1.PodList

		waitErr := wait.PollImmediate(config.Polling.Interval, config.Polling.Timeout, func() (bool, error) {
			pods, err = ctx.Base.Clients.KubeClient.Kube.
				CoreV1().Pods(ns).List(metav1.ListOptions{})
			return inState(pods, err)
		})

		if waitErr != nil {
			t.Error(errors.Wrapf(
				waitErr, "waiting for pods in namespace %s to vanish failed, got: %+v",
				ns, pods))
		}
	})
}
