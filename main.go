package main

import (
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "metrics-logging: ", log.LstdFlags)
	metric := newMetric(logger)
	server := newServer(logger, metric)

	// TODO get address from the config file i.e., environment variable etc.
	listenAddr := ":9000"
	server.start(listenAddr)
}
