package multicast

import (
	"fmt"
	"testing"
	"time"
)

// Additional comprehensive benchmark tests for the multicast package

func BenchmarkMessageOperations(b *testing.B) {
	// Create a sample message for benchmarking
	msg := &Message{
		ID:        12345,
		Timestamp: time.Now(),
		Source:    "benchmark-test-host-name",
	}

	b.Run("Marshal", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := msg.Marshal()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// Pre-marshal for unmarshal tests
	data, err := msg.Marshal()
	if err != nil {
		b.Fatal(err)
	}

	b.Run("Unmarshal", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := UnmarshalMessage(data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Age", func(b *testing.B) {
		// Create message from 1 second ago
		oldMsg := &Message{
			ID:        1,
			Timestamp: time.Now().Add(-time.Second),
			Source:    "test",
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = oldMsg.Age()
		}
	})

	b.Run("FullCycle", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Create message
			msg := &Message{
				ID:        i,
				Timestamp: time.Now(),
				Source:    "benchmark-host",
			}
			
			// Marshal
			data, err := msg.Marshal()
			if err != nil {
				b.Fatal(err)
			}
			
			// Unmarshal
			decoded, err := UnmarshalMessage(data)
			if err != nil {
				b.Fatal(err)
			}
			
			// Calculate age
			_ = decoded.Age()
		}
	})
}

func BenchmarkSenderOperations(b *testing.B) {
	b.Run("NewSender", func(b *testing.B) {
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
	})

	b.Run("NewSenderWithCustomPorts", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sender, err := NewSender("239.23.23.23:2323", "", time.Second, 64, 12345, 8080)
			if err != nil {
				b.Fatal(err)
			}
			if sender != nil {
				sender.conn.Close()
			}
		}
	})

	b.Run("NewSenderWithValidation", func(b *testing.B) {
		// Test the validation overhead
		testCases := []struct {
			ttl   int
			sport int
			dport int
		}{
			{1, 0, 0},
			{255, 65535, 65535},
			{128, 8080, 9090},
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tc := testCases[i%len(testCases)]
			sender, err := NewSender("239.23.23.23:2323", "", time.Second, tc.ttl, tc.sport, tc.dport)
			if err != nil {
				b.Fatal(err)
			}
			if sender != nil {
				sender.conn.Close()
			}
		}
	})
}

func BenchmarkReceiverOperations(b *testing.B) {
	b.Run("NewReceiver", func(b *testing.B) {
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
	})

	b.Run("NewReceiverWithPortOverride", func(b *testing.B) {
		ports := []int{8080, 9090, 7777, 6666}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			port := ports[i%len(ports)]
			receiver, err := NewReceiver("239.23.23.23:2323", "", port)
			if err != nil {
				b.Fatal(err)
			}
			if receiver != nil {
				receiver.conn.Close()
			}
		}
	})
}

func BenchmarkConcurrentOperations(b *testing.B) {
	b.Run("ConcurrentSenderCreation", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sender, err := NewSender("239.23.23.23:2323", "", time.Second, 1, 0, 0)
				if err != nil {
					b.Fatal(err)
				}
				if sender != nil {
					sender.conn.Close()
				}
			}
		})
	})

	b.Run("ConcurrentReceiverCreation", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				receiver, err := NewReceiver("239.23.23.23:2323", "", 0)
				if err != nil {
					b.Fatal(err)
				}
				if receiver != nil {
					receiver.conn.Close()
				}
			}
		})
	})

	b.Run("ConcurrentMessageProcessing", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				msg := &Message{
					ID:        1,
					Timestamp: time.Now(),
					Source:    "concurrent-test",
				}

				data, err := msg.Marshal()
				if err != nil {
					b.Fatal(err)
				}

				decoded, err := UnmarshalMessage(data)
				if err != nil {
					b.Fatal(err)
				}

				_ = decoded.Age()
			}
		})
	})
}

// Memory allocation benchmarks
func BenchmarkMemoryAllocations(b *testing.B) {
	b.Run("MessageCreation", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			_ = &Message{
				ID:        i,
				Timestamp: time.Now(),
				Source:    "allocation-test",
			}
		}
	})

	b.Run("MessageMarshalAllocs", func(b *testing.B) {
		msg := &Message{
			ID:        1,
			Timestamp: time.Now(),
			Source:    "alloc-test",
		}
		
		b.ReportAllocs()
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			_, err := msg.Marshal()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MessageUnmarshalAllocs", func(b *testing.B) {
		msg := &Message{
			ID:        1,
			Timestamp: time.Now(),
			Source:    "alloc-test",
		}
		data, err := msg.Marshal()
		if err != nil {
			b.Fatal(err)
		}
		
		b.ReportAllocs()
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			_, err := UnmarshalMessage(data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Scaling benchmarks with different data sizes
func BenchmarkScaling(b *testing.B) {
	hostnameSizes := []int{10, 50, 100, 500}
	
	for _, size := range hostnameSizes {
		hostname := generateString(size)
		
		b.Run(fmt.Sprintf("Hostname_%d_chars", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				msg := &Message{
					ID:        i,
					Timestamp: time.Now(),
					Source:    hostname,
				}
				
				data, err := msg.Marshal()
				if err != nil {
					b.Fatal(err)
				}
				
				_, err = UnmarshalMessage(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Helper function to generate strings of specific length
func generateString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}