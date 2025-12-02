package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: go run cmd/*.go [command] [--dev]")
		log.Println("Available commands: fetch, fetchAll, update")
		return
	}

	var command string
	devMode := false

	// Simple argument parsing
	for _, arg := range os.Args[1:] {
		if arg == "--dev" {
			devMode = true
			continue
		}
		if command == "" {
			command = arg
		}
	}

	switch command {
	case "fetch":
		runFetch(devMode)
	case "fetchAll":
		runFetchAll(devMode)
	case "update":
		runUpdate(devMode)
	default:
		log.Printf("Unknown command: '%s'. Available commands: fetch, fetchAll, update\n", command)
	}
}
