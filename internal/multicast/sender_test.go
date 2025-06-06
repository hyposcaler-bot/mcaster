package multicast

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSenderValidation(t *testing.T) {
	tests := []struct {
		name      string
		groupAddr string
		iface     string
		interval  time.Duration
		ttl       int
		sport     int
		dport     int
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid parameters",
			groupAddr: "239.23.23.23:2323",
			iface:     "",
			interval:  time.Second,
			ttl:       1,
			sport:     0,
			dport:     0,
			wantErr:   false,
		},
		{
			name:      "valid with custom ports",
			groupAddr: "239.23.23.23:2323",
			iface:     "",
			interval:  time.Second,
			ttl:       64,
			sport:     12345,
			dport:     8080,
			wantErr:   false,
		},
		{
			name:      "invalid TTL low",
			groupAddr: "239.23.23.23:2323",
			iface:     "",
			interval:  time.Second,
			ttl:       0,
			sport:     0,
			dport:     0,
			wantErr:   true,
			errMsg:    "TTL must be between 1 and 255",
		},
		{
			name:      "invalid TTL high",
			groupAddr: "239.23.23.23:2323",
			iface:     "",
			interval:  time.Second,
			ttl:       256,
			sport:     0,
			dport:     0,
			wantErr:   true,
			errMsg:    "TTL must be between 1 and 255",
		},
		{
			name:      "invalid source port negative",
			groupAddr: "239.23.23.23:2323",
			iface:     "",
			interval:  time.Second,
			ttl:       1,
			sport:     -1,
			dport:     0,
			wantErr:   true,
			errMsg:    "source port must be between 0 and 65535",
		},
		{
			name:      "invalid source port high",
			groupAddr: "239.23.23.23:2323",
			iface:     "",
			interval:  time.Second,
			ttl:       1,
			sport:     65536,
			dport:     0,
			wantErr:   true,
			errMsg:    "source port must be between 0 and 65535",
		},
		{
			name:      "invalid destination port high",
			groupAddr: "239.23.23.23:2323",
			iface:     "",
			interval:  time.Second,
			ttl:       1,
			sport:     0,
			dport:     65536,
			wantErr:   true,
			errMsg:    "destination port must be between 1 and 65535",
		},
		{
			name:      "invalid group address",
			groupAddr: "invalid-address",
			iface:     "",
			interval:  time.Second,
			ttl:       1,
			sport:     0,
			dport:     0,
			wantErr:   true,
			errMsg:    "failed to resolve multicast address",
		},
		{
			name:      "invalid address resolution",
			groupAddr: "999.999.999.999:2323",
			iface:     "",
			interval:  time.Second,
			ttl:       1,
			sport:     0,
			dport:     0,
			wantErr:   true,
		},
		{
			name:      "boundary values - valid",
			groupAddr: "239.23.23.23:2323",
			iface:     "",
			interval:  time.Second,
			ttl:       255,
			sport:     65535,
			dport:     65535,
			wantErr:   false,
		},
		{
			name:      "minimum valid TTL",
			groupAddr: "239.23.23.23:2323",
			iface:     "",
			interval:  time.Second,
			ttl:       1,
			sport:     0,
			dport:     0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender, err := NewSender(tt.groupAddr, tt.iface, tt.interval, tt.ttl, tt.sport, tt.dport)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, sender)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, sender)
				
				if sender != nil {
					// Verify sender properties
					assert.Equal(t, tt.ttl, sender.ttl)
					assert.Equal(t, tt.sport, sender.sport)
					assert.Equal(t, tt.interval, sender.interval)
					assert.NotEmpty(t, sender.hostname)
					assert.NotNil(t, sender.conn)
					assert.NotNil(t, sender.groupAddr)
					
					// Clean up
					sender.conn.Close()
				}
			}
		})
	}
}

func TestSenderAddressResolution(t *testing.T) {
	tests := []struct {
		name      string
		groupAddr string
		dport     int
		wantErr   bool
		expectIP  string
		expectPort int
	}{
		{
			name:      "valid IPv4 no override",
			groupAddr: "239.23.23.23:2323",
			dport:     0,
			wantErr:   false,
			expectIP:  "239.23.23.23",
			expectPort: 2323,
		},
		{
			name:      "valid IPv4 with port override",
			groupAddr: "239.23.23.23:2323",
			dport:     8080,
			wantErr:   false,
			expectIP:  "239.23.23.23",
			expectPort: 8080,
		},
		{
			name:      "IPv4 without port, with override",
			groupAddr: "239.23.23.23",
			dport:     9999,
			wantErr:   false,
			expectIP:  "239.23.23.23",
			expectPort: 9999,
		},
		{
			name:      "invalid IPv4 address",
			groupAddr: "300.300.300.300:2323",
			dport:     0,
			wantErr:   true,
		},
		{
			name:      "invalid port in address",
			groupAddr: "239.23.23.23:99999",
			dport:     0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender, err := NewSender(tt.groupAddr, "", time.Second, 1, 0, tt.dport)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, sender)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, sender)
				require.NotNil(t, sender.groupAddr)
				
				assert.Equal(t, tt.expectIP, sender.groupAddr.IP.String())
				assert.Equal(t, tt.expectPort, sender.groupAddr.Port)
				
				// Clean up
				sender.conn.Close()
			}
		})
	}
}

func TestSenderWithInterface(t *testing.T) {
	// Test interface binding errors (most interfaces won't exist in test environment)
	tests := []struct {
		name      string
		iface     string
		wantErr   bool
		errMsg    string
	}{
		{
			name:    "nonexistent interface",
			iface:   "nonexistent-interface-12345",
			wantErr: true,
			errMsg:  "failed to bind to interface",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender, err := NewSender("239.23.23.23:2323", tt.iface, time.Second, 1, 0, 0)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, sender)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, sender)
				if sender != nil {
					sender.conn.Close()
				}
			}
		})
	}
}

func TestSenderStructFields(t *testing.T) {
	// Test that sender struct is properly initialized
	sender, err := NewSender("239.23.23.23:2323", "", 2*time.Second, 32, 12345, 8080)
	require.NoError(t, err)
	require.NotNil(t, sender)
	defer sender.conn.Close()

	// Verify all fields are set correctly
	assert.Equal(t, 32, sender.ttl)
	assert.Equal(t, 12345, sender.sport)
	assert.Equal(t, 2*time.Second, sender.interval)
	assert.Equal(t, 0, sender.packetCount) // Should start at 0
	assert.NotEmpty(t, sender.hostname)
	assert.NotNil(t, sender.conn)
	assert.NotNil(t, sender.groupAddr)
	
	// Verify group address was overridden correctly
	assert.Equal(t, "239.23.23.23", sender.groupAddr.IP.String())
	assert.Equal(t, 8080, sender.groupAddr.Port)
}

func TestSenderHostname(t *testing.T) {
	// Test hostname handling
	sender, err := NewSender("239.23.23.23:2323", "", time.Second, 1, 0, 0)
	require.NoError(t, err)
	require.NotNil(t, sender)
	defer sender.conn.Close()

	// Hostname should be set to something (either actual hostname or "unknown")
	assert.NotEmpty(t, sender.hostname)
	assert.True(t, len(sender.hostname) > 0)
}

func TestSenderConnectionProperties(t *testing.T) {
	// Test that the UDP connection is properly configured
	sender, err := NewSender("239.23.23.23:2323", "", time.Second, 1, 12345, 0)
	require.NoError(t, err)
	require.NotNil(t, sender)
	defer sender.conn.Close()

	// Verify connection properties
	localAddr := sender.conn.LocalAddr()
	assert.NotNil(t, localAddr)

	remoteAddr := sender.conn.RemoteAddr()
	assert.NotNil(t, remoteAddr)
	
	// Verify remote address matches what we set
	udpRemote := remoteAddr.(*net.UDPAddr)
	assert.Equal(t, "239.23.23.23", udpRemote.IP.String())
	assert.Equal(t, 2323, udpRemote.Port)
}

// Test edge cases and error conditions
func TestSenderEdgeCases(t *testing.T) {
	t.Run("zero interval", func(t *testing.T) {
		// Zero interval should be allowed (though not practical)
		sender, err := NewSender("239.23.23.23:2323", "", 0, 1, 0, 0)
		assert.NoError(t, err)
		if sender != nil {
			sender.conn.Close()
		}
	})

	t.Run("very large interval", func(t *testing.T) {
		// Large interval should be allowed
		sender, err := NewSender("239.23.23.23:2323", "", 24*time.Hour, 1, 0, 0)
		assert.NoError(t, err)
		if sender != nil {
			sender.conn.Close()
		}
	})

	t.Run("negative interval", func(t *testing.T) {
		// Negative interval should be allowed (Go time.Ticker handles it)
		sender, err := NewSender("239.23.23.23:2323", "", -time.Second, 1, 0, 0)
		assert.NoError(t, err)
		if sender != nil {
			sender.conn.Close()
		}
	})
}

// Benchmark tests for sender creation
func BenchmarkNewSender(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sender, err := NewSender("239.23.23.23:2323", "", time.Second, 1, 0, 0)
		if err != nil {
			b.Fatal(err)
		}
		if sender != nil {
			sender.conn.Close()
		}
	}
}

func BenchmarkNewSenderWithPortOverride(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sender, err := NewSender("239.23.23.23:2323", "", time.Second, 1, 0, 8080)
		if err != nil {
			b.Fatal(err)
		}
		if sender != nil {
			sender.conn.Close()
		}
	}
}