package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCreateRandomId makes sure the created ticket id is in line with the specification
func TestCreateRandomId(t *testing.T) {

	ticketId := CreateRandomId(10)

	assert.True(t, len(ticketId) == 10, "Random id has the wrong length")
}

