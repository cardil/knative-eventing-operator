package e2e

import (
	"github.com/openshift-knative/knative-eventing-operator/test"
)

// Specifications - specifications that can be executed by serverless operator
// in it's test plan
func Specifications() []test.Specification {
	return specifications
}

var specifications []test.Specification = []test.Specification{
	test.NewSpec("TestKnativeEventingDeployment", testKnativeEventingDeployment),
}
