package multicast

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/yourusername/mcaster/internal/network"
)

// Sender handles multicast packet transmission
type Sender struct {
	conn        *net.UDPConn
	groupAddr   *net.UDPAddr
	hostname    string
	interval    time.Duration
	ttl         int
	sport       int
	packetCount int
}

// NewSender creates a new multicast sender
func NewSender(groupAddr, interfaceName string, interval time.Duration, ttl, sport, dport int) (*Sender, error) {
	// Validate TTL
	if ttl < 1 || ttl > 255 {
		return nil, fmt.Errorf("TTL must be between 1 and 255, got %d", ttl)
	}

	// Validate source port
	if sport < 0 || sport > 65535 {
		return nil, fmt.Errorf("source port must be between 0 and 65535, got %d", sport)
	}

	// Override destination port if specified
	finalGroupAddr, err := network.OverrideGroupPort(groupAddr, dport)
	if err != nil {
		return nil, err
	}

	addr, err := net.ResolveUDPAddr("udp", finalGroupAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve multicast address: %w", err)
	}

	// Create local address with specified source port (0 = random)
	localAddr := &net.UDPAddr{Port: sport}
	conn, err := net.DialUDP("udp", localAddr, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP connection: %w", err)
	}

	// Bind to specific interface if requested
	if interfaceName != "" {
		conn.Close()
		conn, err = network.DialUDPOnInterface(interfaceName, addr, sport)
		if err != nil {
			return nil, fmt.Errorf("failed to bind to interface %s: %w", interfaceName, err)
		}
	}

	// Set TTL for multicast packets
	if err := network.SetMulticastTTL(conn, ttl); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to set multicast TTL: %w", err)
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	return &Sender{
		conn:      conn,
		groupAddr: addr,
		hostname:  hostname,
		interval:  interval,
		ttl:       ttl,
		sport:     sport,
	}, nil
}

// Start begins sending multicast packets
func (s *Sender) Start() error {
	defer s.conn.Close()

	localAddr := s.conn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("🚀 Starting multicast sender to %s\n", s.groupAddr)
	fmt.Printf("📡 Sending packets every %v (TTL: %d, source port: %d)\n", s.interval, s.ttl, localAddr.Port)
	fmt.Printf("⏹️  Press Ctrl+C to stop\n\n")

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := s.sendPacket(); err != nil {
			log.Printf("❌ Failed to send packet: %v", err)
			continue
		}
	}

	return nil
}

func (s *Sender) sendPacket() error {
	s.packetCount++

	msg := Message{
		ID:        s.packetCount,
		Timestamp: time.Now(),
		Source:    s.hostname,
	}

	data, err := msg.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = s.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Printf("📤 [%s] Sent packet #%d\n",
		msg.Timestamp.Format("15:04:05.000"), s.packetCount)

	return nil
}
