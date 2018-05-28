package collector

import (
	"errors"
	"fmt"
	"os/exec"
)

/*
#include <stdio.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <ifaddrs.h>
#include <net/if.h>
*/
import "C"

func getDefaultDev() (map[string]int, error) {
	iface := ""
	netDev := map[string]int{}

	var ifap, ifa *C.struct_ifaddrs
	if C.getifaddrs(&ifap) == -1 {
		return nil, errors.New("getifaddrs() failed")
	}
	defer C.freeifaddrs(ifap)

	for ifa = ifap; ifa != nil; ifa = ifa.ifa_next {
		if ifa.ifa_addr.sa_family == C.AF_LINK {
			dev := C.GoString(ifa.ifa_name)
			netDev[dev] = 0
		}
	}

	out, err := exec.Command("bash", "-c", "route -n get default | grep interface").Output()
	if err != nil {
		return nil, fmt.Errorf("Couldn't get network devices: %s", err)
	}

	fmt.Sscanf(string(out), "  interface:%s", &iface)
	if iface == "" {
		return nil, fmt.Errorf("Parse the interface fail")
	}

	netDev[iface] = 1

	return netDev, nil
}
