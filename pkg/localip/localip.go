package localip

import (
	"fmt"
	"net"
)

func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		addresses, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, address := range addresses {
			ip, ok := address.(*net.IPNet)
			if ok && !ip.IP.IsLoopback() && ip.IP.To4() != nil {
				// Returning the local IP address as a string
				return ip.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no valid local IP address found")
}
