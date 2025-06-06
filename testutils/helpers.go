package testutils

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/hyposcaler-bot/mcaster/internal/multicast"
)

// CreateTestMessage creates a test message with the given ID
func CreateTestMessage(id int) *multicast.Message {
	return &multicast.Message{
		ID:        id,
		Timestamp: time.Now(),
		Source:    "test-host",
	}
}

// CreateTestMessageWithTime creates a test message with specific timestamp
func CreateTestMessageWithTime(id int, timestamp time.Time) *multicast.Message {
	return &multicast.Message{
		ID:        id,
		Timestamp: timestamp,
		Source:    "test-host",
	}
}

// AssertValidMulticastAddr validates that an address is a proper multicast address
func AssertValidMulticastAddr(t *testing.T, addr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	assert.NoError(t, err, "Address should be resolvable")
	assert.True(t, udpAddr.IP.IsMulticast(), "Address should be multicast")
}

// AssertValidPort validates that a port is in the valid range
func AssertValidPort(t *testing.T, port int) {
	assert.GreaterOrEqual(t, port, 0, "Port should be >= 0")
	assert.LessOrEqual(t, port, 65535, "Port should be <= 65535")
}

// GetTestMulticastAddr returns a test multicast address
func GetTestMulticastAddr() string {
	return "239.255.255.255:12345"
}

// GetTestInterface returns the loopback interface for testing
func GetTestInterface() string {
	return "lo"
}