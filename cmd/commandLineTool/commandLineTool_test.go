// Main package of the command line utility
package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
 * Ticketsystem Trivial Tickets
 *
 * Matriculation numbers: 3040018, 6694964, 3478222
 * Lecture:               Programmieren II, INF16B
 * Lecturer:              Herr Prof. Dr. Helmut Neemann
 * Institute:             Duale Hochschule Baden-Württemberg Mosbach
 *
 * ---------------
 *
 * Package main [tests]
 * Main package of the command line utility
 */

func TestGetConfig(t *testing.T) {
	conf, fetch, submit, mail := getConfig()
	assert.NotNil(t, conf)
	assert.Equal(t, "localhost", conf.IPAddr)
	assert.Equal(t, uint16(8443), conf.Port)
	assert.Equal(t, "./ssl/server.cert", conf.Cert)
	assert.False(t, fetch)
	assert.False(t, submit)
	assert.Equal(t, `{"from":"", "subject":"", "message":""}`, mail)
}
