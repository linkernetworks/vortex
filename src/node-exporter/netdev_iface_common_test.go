package collector

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestGetDefaultDev(t *testing.T) {
	iface := ""
	defaultInterface := ""

	switch runtime.GOOS {
	case "linux":
		out, err := exec.Command("bash", "-c", "ip route get 8.8.8.8 | awk 'NR==1 {print $5}'").Output()
		if err != nil {
			t.Error("Can't get the interface for default route, maybe the system does not have the network connectivity", err)
		}
		//the output should look like eth0
		fmt.Sscanf(string(out), "%s", &defaultInterface)
		if defaultInterface == "" {
			t.Error("Parse the interface fail")
		}
	case "darwin":
		// skip the test for osx
		t.Skip("Skip the test when os is drawin.")
	}

	netDevDefault, err := getDefaultDev()
	if err != nil {
		t.Error("couldn't get default network devices", err)
	}

	for key, value := range netDevDefault {
		if value == 1 {
			iface = key
		}
	}

	if strings.Compare(defaultInterface, iface) != 0 {
		t.Error("Default interface didn't match the result from getDefaultDev().")
	}
}
