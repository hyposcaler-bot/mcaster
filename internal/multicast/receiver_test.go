package multicast

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReceiverValidation(t *testing.T) {
	tests := []struct {
		name      string
		groupAddr string
		iface     string
		dport     int
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid parameters",
			groupAddr: "239.23.23.23:2323",
			iface:     "",
			dport:     0,
			wantErr:   false,
		},
		{
			name:      "valid with destination port override",
			groupAddr: "239.23.23.23:2323",
			iface:     "",
			dport:     8080,
			wantErr:   false,
		},
		{
			name:      "valid address without port, with dport",
			groupAddr: "239.23.23.23",
			iface:     "",
			dport:     9999,
			wantErr:   false,
		},
		{
			name:      "invalid destination port high",
			groupAddr: "239.23.23.23:2323",
			iface:     "",
			dport:     65536,
			wantErr:   true,
			errMsg:    "destination port must be between 1 and 65535",
		},
		{
			name:      "invalid group address",
			groupAddr: "invalid-address",
			iface:     "",
			dport:     0,
			wantErr:   true,
			errMsg:    "missing port in address",
		},
		{
			name:      "invalid address resolution",
			groupAddr: "999.999.999.999:2323",
			iface:     "",
			dport:     0,
			wantErr:   true,
		},
		{
			name:      "valid multicast IPv4",
			groupAddr: "224.0.0.1:1234",
			iface:     "",
			dport:     0,
			wantErr:   false,
		},
		{
			name:      "boundary destination port",
			groupAddr: "239.23.23.23:2323",
			iface:     "",
			dport:     65535,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver, err := NewReceiver(tt.groupAddr, tt.iface, tt.dport)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, receiver)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, receiver)
				
				if receiver != nil {
					// Verify receiver properties
					assert.NotNil(t, receiver.conn)
					assert.NotNil(t, receiver.groupAddr)
					assert.NotNil(t, receiver.buffer)
					assert.Equal(t, 1024, len(receiver.buffer)) // Default buffer size
					
					// Clean up
					receiver.conn.Close()
				}
			}
		})
	}
}

func TestReceiverAddressResolution(t *testing.T) {
	tests := []struct {
		name       string
		groupAddr  string
		dport      int
		wantErr    bool
		expectIP   string
		expectPort int
	}{
		{
			name:       "valid IPv4 no override",
			groupAddr:  "239.23.23.23:2323",
			dport:      0,
			wantErr:    false,
			expectIP:   "239.23.23.23",
			expectPort: 2323,
		},
		{
			name:       "valid IPv4 with port override",
			groupAddr:  "239.23.23.23:2323",
			dport:      8080,
			wantErr:    false,
			expectIP:   "239.23.23.23",
			expectPort: 8080,
		},
		{
			name:       "IPv4 without port, with override",
			groupAddr:  "239.23.23.23",
			dport:      9999,
			wantErr:    false,
			expectIP:   "239.23.23.23",
			expectPort: 9999,
		},
		{
			name:      "invalid IPv4 address",
			groupAddr: "300.300.300.300:2323",
			dport:     0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver, err := NewReceiver(tt.groupAddr, "", tt.dport)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, receiver)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, receiver)
				require.NotNil(t, receiver.groupAddr)
				
				assert.Equal(t, tt.expectIP, receiver.groupAddr.IP.String())
				assert.Equal(t, tt.expectPort, receiver.groupAddr.Port)
				
				// Verify it's a multicast address
				assert.True(t, receiver.groupAddr.IP.IsMulticast())
				
				// Clean up
				receiver.conn.Close()
			}
		})
	}
}

func TestReceiverWithInterface(t *testing.T) {
	// Test interface binding (most interfaces won't exist in test environment)
	tests := []struct {
		name    string
		iface   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nonexistent interface",
			iface:   "nonexistent-interface-12345",
			wantErr: true,
			errMsg:  "failed to find interface",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver, err := NewReceiver("239.23.23.23:2323", tt.iface, 0)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, receiver)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, receiver)
				if receiver != nil {
					receiver.conn.Close()
				}
			}
		})
	}
}

func TestReceiverStructFields(t *testing.T) {
	// Test that receiver struct is properly initialized
	receiver, err := NewReceiver("239.23.23.23:8080", "", 0)
	require.NoError(t, err)
	require.NotNil(t, receiver)
	defer receiver.conn.Close()

	// Verify all fields are set correctly
	assert.NotNil(t, receiver.conn)
	assert.NotNil(t, receiver.groupAddr)
	assert.NotNil(t, receiver.buffer)
	assert.Equal(t, 1024, len(receiver.buffer))
	
	// Verify group address
	assert.Equal(t, "239.23.23.23", receiver.groupAddr.IP.String())
	assert.Equal(t, 8080, receiver.groupAddr.Port)
}

func TestReceiverBufferSize(t *testing.T) {
	// Test that the buffer is allocated with correct size
	receiver, err := NewReceiver("239.23.23.23:2323", "", 0)
	require.NoError(t, err)
	require.NotNil(t, receiver)
	defer receiver.conn.Close()

	// Verify buffer properties
	assert.NotNil(t, receiver.buffer)
	assert.Equal(t, 1024, len(receiver.buffer))
	assert.Equal(t, 1024, cap(receiver.buffer))
	
	// Buffer should be zero-initialized
	for i, b := range receiver.buffer {
		if b != 0 {
			t.Errorf("Buffer byte at index %d should be 0, got %d", i, b)
			break
		}
	}
}

func TestReceiverConnectionProperties(t *testing.T) {
	// Test that the UDP connection is properly configured for multicast
	receiver, err := NewReceiver("239.23.23.23:2323", "", 0)
	require.NoError(t, err)
	require.NotNil(t, receiver)
	defer receiver.conn.Close()

	// Verify connection properties
	localAddr := receiver.conn.LocalAddr()
	assert.NotNil(t, localAddr)
	
	// For multicast receiver, local address should be listening on the multicast port
	udpLocal := localAddr.(*net.UDPAddr)
	assert.Equal(t, 2323, udpLocal.Port)
}

func TestReceiverMulticastAddressTypes(t *testing.T) {
	// Test different types of multicast addresses
	multicastAddresses := []struct {
		name string
		addr string
		desc string
	}{
		{
			name: "local network control",
			addr: "224.0.0.1:2323",
			desc: "All Systems multicast",
		},
		{
			name: "internetwork control",
			addr: "224.0.1.1:2323", 
			desc: "Internetwork control block",
		},
		{
			name: "organization local",
			addr: "239.255.255.255:2323",
			desc: "Organization-local scope",
		},
		{
			name: "default test address",
			addr: "239.23.23.23:2323",
			desc: "Our default test address",
		},
	}

	for _, addr := range multicastAddresses {
		t.Run(addr.name, func(t *testing.T) {
			receiver, err := NewReceiver(addr.addr, "", 0)
			
			// Note: Some of these might fail depending on system configuration
			// This test documents expected behavior
			if err != nil {
				t.Logf("Expected potential failure for %s (%s): %v", addr.desc, addr.addr, err)
				return
			}
			
			require.NotNil(t, receiver)
			assert.True(t, receiver.groupAddr.IP.IsMulticast())
			receiver.conn.Close()
		})
	}
}

// Test edge cases and error conditions
func TestReceiverEdgeCases(t *testing.T) {
	t.Run("IPv6 multicast address", func(t *testing.T) {
		// Test IPv6 multicast (might not be supported in all environments)
		receiver, err := NewReceiver("[ff02::1]:2323", "", 0)
		if err != nil {
			// IPv6 might not be available, that's okay
			t.Logf("IPv6 multicast not available: %v", err)
			return
		}
		
		if receiver != nil {
			assert.True(t, receiver.groupAddr.IP.IsMulticast())
			receiver.conn.Close()
		}
	})

	t.Run("port zero in address", func(t *testing.T) {
		// Test behavior with port 0 (multicast might use port 0)
		receiver, err := NewReceiver("239.23.23.23:0", "", 0)
		if err != nil {
			// Port 0 might not be allowed for multicast
			t.Logf("Port 0 not allowed: %v", err)
			return
		}
		
		if receiver != nil {
			// Port 0 is actually allowed for multicast addresses
			assert.Equal(t, 0, receiver.groupAddr.Port)
			receiver.conn.Close()
		}
	})
}

// Integration-style tests
func TestReceiverIntegration(t *testing.T) {
	t.Run("address override and resolution", func(t *testing.T) {
		// Test the complete flow of overriding port and creating receiver
		originalAddr := "239.23.23.23:2323"
		newPort := 8080

		receiver, err := NewReceiver(originalAddr, "", newPort)
		require.NoError(t, err)
		require.NotNil(t, receiver)
		defer receiver.conn.Close()

		// Verify the port was overridden
		assert.Equal(t, "239.23.23.23", receiver.groupAddr.IP.String())
		assert.Equal(t, newPort, receiver.groupAddr.Port)
		assert.True(t, receiver.groupAddr.IP.IsMulticast())
	})
}

// Benchmark tests for receiver creation
func BenchmarkNewReceiver(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		receiver, err := NewReceiver("239.23.23.23:2323", "", 0)
		if err != nil {
			b.Fatal(err)
		}
		if receiver != nil {
			receiver.conn.Close()
		}
	}
}

func BenchmarkNewReceiverWithPortOverride(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		receiver, err := NewReceiver("239.23.23.23:2323", "", 8080)
		if err != nil {
			b.Fatal(err)
		}
		if receiver != nil {
			receiver.conn.Close()
		}
	}
}