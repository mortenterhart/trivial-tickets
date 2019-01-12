package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
