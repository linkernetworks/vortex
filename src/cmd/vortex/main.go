package main

import (
	"bitbucket.org/linkernetworks/vortex/src/server"
	"flag"
)

func main() {
	var (
		configPath string
		host       string
		port       string
	)

	flag.StringVar(&configPath, "config", "config/local.json", "config file path")
	flag.StringVar(&host, "host", "0.0.0.0", "hostname")
	flag.StringVar(&port, "port", "7890", "port")

	a := server.App{}
	a.LoadConfig(configPath).Start(host, port)
}
