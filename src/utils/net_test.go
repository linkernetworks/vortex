package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPToCIDR(t *testing.T) {
	ip := "1.2.3.4"
	netmask := "255.255.240.0"
	c := IPToCIDR(ip, netmask)
	assert.Equal(t, c, "1.2.3.4/20")
}
