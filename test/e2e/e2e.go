package e2e

import (
	"testing"

	"github.com/openshift-knative/knative-eventing-operator/test"
	pkgTest "knative.dev/pkg/test"
)

// Setup creates the client objects needed in the e2e tests.
func Setup(t *testing.T) *test.Clients {
	clients, err := test.NewClients(
		pkgTest.Flags.Kubeconfig,
		pkgTest.Flags.Cluster,
	)
	if err != nil {
		t.Fatalf("Couldn't initialize clients: %v", err)
	}
	return clients
}
