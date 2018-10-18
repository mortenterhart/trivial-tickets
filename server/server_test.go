package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

func TestGetTemplates(t *testing.T) {

	tmpl := GetTemplates("../www")

	assert.NotNil(t, tmpl, "GetTemplates() returned no found templates")
}
