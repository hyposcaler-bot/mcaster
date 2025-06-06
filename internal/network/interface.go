package network

import (
	"fmt"
	"log"
	"net"
)

// GetInterface returns a network interface by name, or nil for default
func GetInterface(interfaceName string) *net.Interface {
	if interfaceName == "" {
		return nil
	}

	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		log.Fatalf("Failed to find interface %s: %v", interfaceName, err)
	}

	return iface
}

// DialUDPOnInterface creates a UDP connection bound to a specific interface
func DialUDPOnInterface(interfaceName string, remoteAddr *net.UDPAddr) (*net.UDPConn, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("failed to find interface %s: %w", interfaceName, err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("failed to get interface addresses: %w", err)
	}

	var localAddr *net.UDPAddr
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localAddr = &net.UDPAddr{IP: ipnet.IP, Port: 0}
				break
			}
		}
	}

	if localAddr == nil {
		return nil, fmt.Errorf("no suitable IPv4 address found on interface %s", interfaceName)
	}

	return net.DialUDP("udp", localAddr, remoteAddr)
}
