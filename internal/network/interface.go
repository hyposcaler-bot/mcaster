package network

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"syscall"
)

// GetInterface returns a network interface by name, or nil for default
func GetInterface(interfaceName string) (*net.Interface, error) {
	if interfaceName == "" {
		return nil, nil
	}

	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("failed to find interface %s: %w", interfaceName, err)
	}

	return iface, nil
}

// DialUDPOnInterface creates a UDP connection bound to a specific interface
func DialUDPOnInterface(interfaceName string, remoteAddr *net.UDPAddr, sport int) (*net.UDPConn, error) {
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
				localAddr = &net.UDPAddr{IP: ipnet.IP, Port: sport}
				break
			}
		}
	}

	if localAddr == nil {
		return nil, fmt.Errorf("no suitable IPv4 address found on interface %s", interfaceName)
	}

	return net.DialUDP("udp", localAddr, remoteAddr)
}

// SetMulticastTTL sets the TTL for multicast packets on a UDP connection
func SetMulticastTTL(conn *net.UDPConn, ttl int) error {
	rawConn, err := conn.SyscallConn()
	if err != nil {
		return fmt.Errorf("failed to get raw connection: %w", err)
	}

	var setErr error
	err = rawConn.Control(func(fd uintptr) {
		// Set IP_MULTICAST_TTL socket option
		setErr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_MULTICAST_TTL, ttl)
	})

	if err != nil {
		return fmt.Errorf("failed to control socket: %w", err)
	}
	if setErr != nil {
		return fmt.Errorf("failed to set multicast TTL: %w", setErr)
	}

	return nil
}

// OverrideGroupPort overrides the port in a group address string if dport > 0
func OverrideGroupPort(groupAddr string, dport int) (string, error) {
	if dport <= 0 {
		return groupAddr, nil
	}

	if dport > 65535 {
		return "", fmt.Errorf("destination port must be between 1 and 65535, got %d", dport)
	}

	// Parse the address to get the host part
	host, _, err := net.SplitHostPort(groupAddr)
	if err != nil {
		// If there's no port, assume it's just the host
		if strings.Contains(err.Error(), "missing port") {
			return net.JoinHostPort(groupAddr, strconv.Itoa(dport)), nil
		}
		// Handle IPv6 addresses without brackets that have parsing issues
		if strings.Contains(err.Error(), "too many colons") {
			return net.JoinHostPort(groupAddr, strconv.Itoa(dport)), nil
		}
		return "", fmt.Errorf("failed to parse group address: %w", err)
	}

	// Replace with the new port
	return net.JoinHostPort(host, strconv.Itoa(dport)), nil
}
