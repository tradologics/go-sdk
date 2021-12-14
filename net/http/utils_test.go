package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var fun = includeProtocol

func TestIncludeProtocol(t *testing.T) {
	httpUrl := "http://example.com"
	httpsUrl := "https://example.com"
	tcpUrl := "tcp://example.com"
	invalidUrl := "example.com"

	msg := "invalid protocol"
	assert.Equal(t, true, fun(httpUrl), msg)
	assert.Equal(t, true, fun(httpsUrl), msg)
	assert.Equal(t, true, fun(tcpUrl), msg)
	assert.Equal(t, false, fun(invalidUrl), msg)

}

func TestDoNotIncludeProtocol(t *testing.T) {
	exp1 := "http:/example.com"
	exp2 := "https//example.com"
	exp3 := "tcp:\\example.com"
	validUrl := "https://example.com"

	msg := "contains protocol"
	assert.Equal(t, false, fun(exp1), msg)
	assert.Equal(t, false, fun(exp2), msg)
	assert.Equal(t, false, fun(exp3), msg)
	assert.Equal(t, true, fun(validUrl), msg)
}
