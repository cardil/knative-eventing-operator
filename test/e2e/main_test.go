//+build e2e

package e2e

import (
	"os"
	"testing"

	"github.com/cardil/operator-e2e-harness/pkg/catalogsource"
)

// TestMain is a main test runner
func TestMain(m *testing.M) {
	catalogsource.Deploy()
	code := m.Run()
	catalogsource.Undeploy()
	os.Exit(code)
}
