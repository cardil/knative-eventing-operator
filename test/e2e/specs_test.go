package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpec(t *testing.T) {
	specs := Specifications()

	assert.NotEmpty(t, specs)
}
