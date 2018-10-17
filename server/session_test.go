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

func TestCreateSession(t *testing.T) {

	session := CreateSession()

	assert.NotNil(t, session, "Session is nil")
}

func TestDestroySession(t *testing.T) {

}

func TestIsSessionValid(t *testing.T) {

}
