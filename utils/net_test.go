package utils_test

import (
	"testing"

	"github.com/shatteredsilicon/ssm-managed/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsTCPAddressLocal(t *testing.T) {
	addrs := map[string]bool{
		"127.0.0.1:80":   true,
		"8.8.8.8:8080":   false,
		"localhost:3306": true,
		"google.com:80":  false,
	}
	for addr, expected := range addrs {
		assert.Equal(t, expected, utils.IsTCPAddressLocal(addr), "addr = %s", addr)
	}
}
