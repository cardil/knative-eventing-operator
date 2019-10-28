// Code generated by injection-gen. DO NOT EDIT.

package fake

import (
	"context"

	knativeeventing "github.com/openshift-knative/knative-eventing-operator/pkg/client/injection/informers/eventing/v1alpha1/knativeeventing"
	fake "github.com/openshift-knative/knative-eventing-operator/pkg/client/injection/informers/factory/fake"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
)

var Get = knativeeventing.Get

func init() {
	injection.Fake.RegisterInformer(withInformer)
}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := fake.Get(ctx)
	inf := f.Eventing().V1alpha1().KnativeEventings()
	return context.WithValue(ctx, knativeeventing.Key{}, inf), inf.Informer()
}
