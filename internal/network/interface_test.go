package network

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOverrideGroupPort(t *testing.T) {
	tests := []struct {
		name      string
		groupAddr string
		dport     int
		expected  string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "no override (dport = 0)",
			groupAddr: "239.23.23.23:2323",
			dport:     0,
			expected:  "239.23.23.23:2323",
			wantErr:   false,
		},
		{
			name:      "override with existing port",
			groupAddr: "239.23.23.23:2323",
			dport:     8080,
			expected:  "239.23.23.23:8080",
			wantErr:   false,
		},
		{
			name:      "override address without port",
			groupAddr: "239.23.23.23",
			dport:     8080,
			expected:  "239.23.23.23:8080",
			wantErr:   false,
		},
		{
			name:      "negative dport is ignored",
			groupAddr: "239.23.23.23:2323",
			dport:     -1,
			expected:  "239.23.23.23:2323",
			wantErr:   false,
		},
		{
			name:      "invalid high port",
			groupAddr: "239.23.23.23:2323",
			dport:     65536,
			expected:  "",
			wantErr:   true,
			errMsg:    "destination port must be between 1 and 65535",
		},
		{
			name:      "IPv6 address with port",
			groupAddr: "[ff02::1]:2323",
			dport:     8080,
			expected:  "[ff02::1]:8080",
			wantErr:   false,
		},
		{
			name:      "IPv6 address without port",
			groupAddr: "ff02::1",
			dport:     8080,
			expected:  "[ff02::1]:8080",
			wantErr:   false,
		},
		{
			name:      "valid port boundaries",
			groupAddr: "239.23.23.23:2323",
			dport:     65535,
			expected:  "239.23.23.23:65535",
			wantErr:   false,
		},
		{
			name:      "port 1",
			groupAddr: "239.23.23.23:2323",
			dport:     1,
			expected:  "239.23.23.23:1",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := OverrideGroupPort(tt.groupAddr, tt.dport)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" && err != nil {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)

				// Verify the result is a valid UDP address
				if result != "" {
					_, resolveErr := net.ResolveUDPAddr("udp", result)
					assert.NoError(t, resolveErr, "Result should be a valid UDP address")
				}
			}
		})
	}
}

func TestGetInterface(t *testing.T) {
	tests := []struct {
		name          string
		interfaceName string
		expectNil     bool
		expectError   bool
	}{
		{
			name:          "empty interface name returns nil",
			interfaceName: "",
			expectNil:     true,
			expectError:   false,
		},
		{
			name:          "loopback interface exists",
			interfaceName: "lo",
			expectNil:     false,
			expectError:   false,
		},
		{
			name:          "nonexistent interface returns error",
			interfaceName: "nonexistent-interface-12345",
			expectNil:     true,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetInterface(tt.interfaceName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectNil {
					assert.Nil(t, result)
				} else {
					// Note: This test might fail on systems without 'lo' interface
					// In real testing environments, we'd need to discover available interfaces
					if result != nil {
						assert.Equal(t, tt.interfaceName, result.Name)
					}
				}
			}
		})
	}
}

func TestDialUDPOnInterface(t *testing.T) {
	// This test is more complex as it requires actual network interfaces
	// We'll test the error cases that don't require real network operations

	tests := []struct {
		name          string
		interfaceName string
		remoteAddr    string
		sport         int
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "invalid interface name",
			interfaceName: "nonexistent-interface-12345",
			remoteAddr:    "239.23.23.23:2323",
			sport:         0,
			wantErr:       true,
			errMsg:        "failed to find interface",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := net.ResolveUDPAddr("udp", tt.remoteAddr)
			require.NoError(t, err)

			conn, err := DialUDPOnInterface(tt.interfaceName, addr, tt.sport)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, conn)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, conn)
				if conn != nil {
					conn.Close()
				}
			}
		})
	}
}

// Integration test for network functionality
func TestNetworkIntegration(t *testing.T) {
	t.Run("address override and resolution", func(t *testing.T) {
		// Test the complete flow of overriding port and resolving address
		originalAddr := "239.23.23.23:2323"
		newPort := 8080

		overriddenAddr, err := OverrideGroupPort(originalAddr, newPort)
		require.NoError(t, err)
		assert.Equal(t, "239.23.23.23:8080", overriddenAddr)

		// Verify it resolves to a valid UDP address
		udpAddr, err := net.ResolveUDPAddr("udp", overriddenAddr)
		require.NoError(t, err)
		assert.True(t, udpAddr.IP.IsMulticast())
		assert.Equal(t, newPort, udpAddr.Port)
	})
}

// Benchmark tests for performance-critical functions
func BenchmarkOverrideGroupPort(b *testing.B) {
	groupAddr := "239.23.23.23:2323"
	dport := 8080

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := OverrideGroupPort(groupAddr, dport)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOverrideGroupPortNoChange(b *testing.B) {
	groupAddr := "239.23.23.23:2323"
	dport := 0 // No override

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := OverrideGroupPort(groupAddr, dport)
		if err != nil {
			b.Fatal(err)
		}
	}
}