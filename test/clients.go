package test

import (
	harnesscfg "github.com/cardil/operator-e2e-harness/pkg/config"
	harness "github.com/cardil/operator-e2e-harness/pkg/test"
	"testing"
)

func init() {
	harnesscfg.Channel = "alpha"
}

// EventingClients contains a clients specific to eventing
type EventingClients struct {
	// TODO: add eventing client to simplify e2e tests
}

// EventingContext holds a context of eventing e2e test
type EventingContext struct {
	Base    *harness.Context
	Clients *EventingClients
}

// NewContext creates a new eventing context
func NewContext(t *testing.T) *EventingContext {
	ctx := harness.NewContext(t, "Admin", 0)
	return &EventingContext{
		Base: ctx,
		Clients: &EventingClients{},
	}
}

// Cleanup invokes a base type cleanup operations
func (ctx *EventingContext) Cleanup() {
	ctx.Base.Cleanup()
}
