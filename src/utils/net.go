package utils

import (
	"fmt"
	"net"
)

// IPToCIDR will do like 0.0.0.0/255.255.255.0 to 0.0.0.0/24
func IPToCIDR(ip string, netmask string) string {
	mask := net.IPMask(net.ParseIP(netmask).To4())
	size, _ := mask.Size()
	return fmt.Sprintf("%s/%d", ip, size)
}
