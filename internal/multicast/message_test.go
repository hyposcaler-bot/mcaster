package multicast

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageMarshalUnmarshal(t *testing.T) {
	now := time.Now()
	msg := &Message{
		ID:        42,
		Timestamp: now,
		Source:    "test-host",
	}

	// Test Marshal
	data, err := msg.Marshal()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Verify it's valid JSON
	var temp map[string]interface{}
	err = json.Unmarshal(data, &temp)
	require.NoError(t, err)

	// Test Unmarshal
	decoded, err := UnmarshalMessage(data)
	require.NoError(t, err)
	require.NotNil(t, decoded)

	// Verify all fields
	assert.Equal(t, msg.ID, decoded.ID)
	assert.Equal(t, msg.Source, decoded.Source)
	
	// Time comparison with tolerance for JSON serialization precision
	assert.WithinDuration(t, msg.Timestamp, decoded.Timestamp, time.Millisecond)
}

func TestMessageAge(t *testing.T) {
	tests := []struct {
		name      string
		timeAgo   time.Duration
		tolerance time.Duration
	}{
		{
			name:      "5 seconds ago",
			timeAgo:   5 * time.Second,
			tolerance: 100 * time.Millisecond,
		},
		{
			name:      "1 minute ago",
			timeAgo:   time.Minute,
			tolerance: 100 * time.Millisecond,
		},
		{
			name:      "just now",
			timeAgo:   0,
			tolerance: 10 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			past := time.Now().Add(-tt.timeAgo)
			msg := &Message{
				ID:        1,
				Timestamp: past,
				Source:    "test-host",
			}

			age := msg.Age()
			
			// Age should be approximately the expected duration
			assert.InDelta(t, tt.timeAgo.Seconds(), age.Seconds(), tt.tolerance.Seconds())
			assert.True(t, age >= 0, "Age should not be negative")
		})
	}
}

func TestUnmarshalInvalidJSON(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
	}{
		{
			name:        "invalid JSON syntax",
			data:        []byte("invalid json"),
			expectError: true,
		},
		{
			name:        "empty data",
			data:        []byte(""),
			expectError: true,
		},
		{
			name:        "null data",
			data:        nil,
			expectError: true,
		},
		{
			name:        "incomplete JSON",
			data:        []byte(`{"id": 1, "timestamp":`),
			expectError: true,
		},
		{
			name:        "wrong JSON structure",
			data:        []byte(`{"wrong": "fields"}`),
			expectError: false, // JSON unmarshaling may succeed but with zero values
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := UnmarshalMessage(tt.data)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, msg)
			} else {
				// For cases where JSON is valid but has wrong structure
				if err == nil {
					assert.NotNil(t, msg)
					// Zero values are expected for wrong structure
					assert.Equal(t, 0, msg.ID)
					assert.Empty(t, msg.Source)
				} else {
					assert.Error(t, err)
					assert.Nil(t, msg)
				}
			}
		})
	}
}

func TestMessageWithDifferentTypes(t *testing.T) {
	tests := []struct {
		name   string
		id     int
		source string
	}{
		{
			name:   "zero ID",
			id:     0,
			source: "host-zero",
		},
		{
			name:   "negative ID",
			id:     -1,
			source: "host-negative",
		},
		{
			name:   "large ID",
			id:     999999,
			source: "host-large",
		},
		{
			name:   "empty source",
			id:     1,
			source: "",
		},
		{
			name:   "unicode source",
			id:     1,
			source: "host-测试-🚀",
		},
		{
			name:   "long source",
			id:     1,
			source: "very-long-hostname-that-might-cause-issues-in-serialization-process",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{
				ID:        tt.id,
				Timestamp: time.Now(),
				Source:    tt.source,
			}

			data, err := msg.Marshal()
			require.NoError(t, err)

			decoded, err := UnmarshalMessage(data)
			require.NoError(t, err)

			assert.Equal(t, tt.id, decoded.ID)
			assert.Equal(t, tt.source, decoded.Source)
		})
	}
}

func TestMessageTimestampPrecision(t *testing.T) {
	// Test various timestamp precisions
	timestamps := []time.Time{
		time.Now(),
		time.Now().Round(time.Second),
		time.Now().Round(time.Millisecond),
		time.Now().Round(time.Microsecond),
		time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC),
	}

	for i, ts := range timestamps {
		t.Run(ts.String(), func(t *testing.T) {
			msg := &Message{
				ID:        i,
				Timestamp: ts,
				Source:    "test-host",
			}

			data, err := msg.Marshal()
			require.NoError(t, err)

			decoded, err := UnmarshalMessage(data)
			require.NoError(t, err)

			// JSON timestamp precision might vary, so allow some tolerance
			assert.WithinDuration(t, ts, decoded.Timestamp, time.Millisecond)
		})
	}
}

func TestMessageJSONFormat(t *testing.T) {
	// Test that the JSON format matches expectations
	msg := &Message{
		ID:        123,
		Timestamp: time.Date(2023, 6, 15, 14, 30, 45, 0, time.UTC),
		Source:    "test-host",
	}

	data, err := msg.Marshal()
	require.NoError(t, err)

	// Parse JSON to verify structure
	var parsed map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	// Verify expected fields exist
	assert.Contains(t, parsed, "id")
	assert.Contains(t, parsed, "timestamp")
	assert.Contains(t, parsed, "source")

	// Verify field types
	assert.IsType(t, float64(0), parsed["id"]) // JSON numbers are float64
	assert.IsType(t, "", parsed["timestamp"])
	assert.IsType(t, "", parsed["source"])

	// Verify values
	assert.Equal(t, float64(123), parsed["id"])
	assert.Equal(t, "test-host", parsed["source"])
}

// Benchmark tests for performance-critical message operations
func BenchmarkMessageMarshal(b *testing.B) {
	msg := &Message{
		ID:        1,
		Timestamp: time.Now(),
		Source:    "benchmark-host",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := msg.Marshal()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMessageUnmarshal(b *testing.B) {
	msg := &Message{
		ID:        1,
		Timestamp: time.Now(),
		Source:    "benchmark-host",
	}
	data, err := msg.Marshal()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := UnmarshalMessage(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMessageAge(b *testing.B) {
	msg := &Message{
		ID:        1,
		Timestamp: time.Now().Add(-5 * time.Second),
		Source:    "benchmark-host",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = msg.Age()
	}
}

func BenchmarkMessageMarshalUnmarshal(b *testing.B) {
	msg := &Message{
		ID:        1,
		Timestamp: time.Now(),
		Source:    "benchmark-host",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, err := msg.Marshal()
		if err != nil {
			b.Fatal(err)
		}
		
		_, err = UnmarshalMessage(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}