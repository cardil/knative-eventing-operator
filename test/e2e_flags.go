// This file contains logic to encapsulate flags which are needed to specify
// what cluster, etc. to use for e2e tests.

package test

const (
	// EventingOperatorNamespace is the default namespace for eventing operator e2e tests
	EventingOperatorNamespace = "operator-tests"
	// EventingOperatorName is the default operator name for eventing operator e2e tests
	EventingOperatorName = "knative-eventing"
)
