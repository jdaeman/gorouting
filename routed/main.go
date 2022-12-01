package main

import (
	"flag"
	"modules/server"
	"os"
)

func parseCommand() *server.ServerConfig {

	dataDir := flag.String("d", "", "data directory")
	port := flag.Int("p", 8000, "server port")
	help := flag.Bool("h", false, "Help message")

	flag.Parse()
	if *dataDir == "" || *help {
		flag.PrintDefaults()
		return nil
	}

	config := &server.ServerConfig{
		DataPath: *dataDir,
		Port:     int16(*port),
	}

	return config
}

func main() {
	config := parseCommand()
	if config == nil {
		os.Exit(2)
	}
	server.Run(*config)
}
