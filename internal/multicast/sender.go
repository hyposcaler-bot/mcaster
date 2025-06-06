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
	packetCount int
}

// NewSender creates a new multicast sender
func NewSender(groupAddr, interfaceName string, interval time.Duration) (*Sender, error) {
	addr, err := net.ResolveUDPAddr("udp", groupAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve multicast address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP connection: %w", err)
	}

	// Bind to specific interface if requested
	if interfaceName != "" {
		conn.Close()
		conn, err = network.DialUDPOnInterface(interfaceName, addr)
		if err != nil {
			return nil, fmt.Errorf("failed to bind to interface %s: %w", interfaceName, err)
		}
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
	}, nil
}

// Start begins sending multicast packets
func (s *Sender) Start() error {
	defer s.conn.Close()

	fmt.Printf("🚀 Starting multicast sender to %s\n", s.groupAddr)
	fmt.Printf("📡 Sending packets every %v\n", s.interval)
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
