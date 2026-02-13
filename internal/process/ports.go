package process

import (
	"fmt"
	"net"
)

const (
	DefaultPort  = 3456
	PortRangeEnd = 3500
)

func IsPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func FindAvailablePort(startPort int) (int, error) {
	for port := startPort; port <= PortRangeEnd; port++ {
		if IsPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port in range %d-%d", startPort, PortRangeEnd)
}
