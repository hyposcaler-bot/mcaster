package multicast

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/yourusername/mcaster/internal/network"
)

// Receiver handles multicast packet reception
type Receiver struct {
	conn      *net.UDPConn
	groupAddr *net.UDPAddr
	buffer    []byte
}

// NewReceiver creates a new multicast receiver
func NewReceiver(groupAddr, interfaceName string, dport int) (*Receiver, error) {
	// Override destination port if specified
	finalGroupAddr, err := network.OverrideGroupPort(groupAddr, dport)
	if err != nil {
		return nil, err
	}

	addr, err := net.ResolveUDPAddr("udp", finalGroupAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve multicast address: %w", err)
	}

	iface, err := network.GetInterface(interfaceName)
	if err != nil {
		return nil, err
	}
	
	conn, err := net.ListenMulticastUDP("udp", iface, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on multicast address: %w", err)
	}

	return &Receiver{
		conn:      conn,
		groupAddr: addr,
		buffer:    make([]byte, 1024),
	}, nil
}

// Start begins receiving multicast packets
func (r *Receiver) Start() error {
	defer r.conn.Close()

	fmt.Printf("🎯 Starting multicast receiver on %s\n", r.groupAddr)
	fmt.Printf("👂 Waiting for packets...\n\n")

	for {
		if err := r.receivePacket(); err != nil {
			log.Printf("❌ Failed to receive packet: %v", err)
			continue
		}
	}
}

func (r *Receiver) receivePacket() error {
	n, remoteAddr, err := r.conn.ReadFromUDP(r.buffer)
	if err != nil {
		return fmt.Errorf("failed to read UDP message: %w", err)
	}

	msg, err := UnmarshalMessage(r.buffer[:n])
	if err != nil {
		fmt.Printf("📥 [%s] Received %d bytes from %s (invalid JSON): %s\n",
			time.Now().Format("15:04:05.000"), n, remoteAddr, string(r.buffer[:n]))
		return nil
	}

	fmt.Printf("📥 [%s] Received packet #%d from %s (%s) - delay: %v\n",
		time.Now().Format("15:04:05.000"), msg.ID, msg.Source, remoteAddr, msg.Age())

	return nil
}
