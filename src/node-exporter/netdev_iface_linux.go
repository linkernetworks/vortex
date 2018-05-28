package collector

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func getDefaultDev() (map[string]int, error) {
	// /proc/net/route stores the kernel's routing table
	// The interface whose destination is 00000000 is the interface of the default gateway
	file, err := os.Open("/proc/net/route")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	netDev, err := getNetDev()
	if err != nil {
		return nil, fmt.Errorf("couldn't get network devices: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	for scanner.Scan() {
		if scanner.Text() == "" {
			break
		}
		s := strings.FieldsFunc(scanner.Text(), Split)
		if s[1] == "00000000" {
			netDev[s[0]] = 1
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return netDev, nil

}
