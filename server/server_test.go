package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTemplates(t *testing.T) {

	tmpl := GetTemplates()

	assert.NotNil(t, tmpl, "GetTemplates() returned no found templates")
}
