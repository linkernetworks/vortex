//phstsai
package collector

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type netDevCollectorTest struct {
	subsystem   string
	metricDescs map[string]*prometheus.Desc
}

func init() {
	registerCollector("iface", defaultEnabled, NewNetDevCollectorTest)
}

// NewNetDevCollectorTest returns a new Collector exposing network device stats.
func NewNetDevCollectorTest() (Collector, error) {
	return &netDevCollectorTest{
		subsystem:   "network",
		metricDescs: map[string]*prometheus.Desc{},
	}, nil
}

func (c *netDevCollectorTest) Update(ch chan<- prometheus.Metric) error {
	netDevDefault, err := getDefaultDev()
	if err != nil {
		return fmt.Errorf("couldn't get default network devices: %s", err)
	}

	for key, value := range netDevDefault {
		desc, ok := c.metricDescs[key]

		if !ok {
			desc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, c.subsystem, "interface"),
				fmt.Sprintf("List network devices and label the default one."),
				[]string{"device"},
				nil,
			)
			c.metricDescs[key] = desc
		}

		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(value), key)
	}
	return nil
}

func Split(r rune) bool {
	return r == ' ' || r == '\t'
}

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
	}

	return netDev, nil
}

func getNetDev() (map[string]int, error) {
	netDev := map[string]int{}

	ifaces, _ := net.Interfaces()

	for _, iface := range ifaces {
		netDev[iface.Name] = 0
	}

	return netDev, nil
}
