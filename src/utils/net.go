package utils

import (
	"fmt"
	"net"
)

func IPToCIDR(ip string, netmask string) string {
	mask := net.IPMask(net.ParseIP(netmask).To4())
	size, _ := mask.Size()
	return fmt.Sprintf("%s/%d", ip, size)
}
