// Package main is the entry point for the GophKeeper application.
package main

import (
	"flag"
	"fmt"
	"os"

	clientCmd "github.com/Pklerik/gophKeep/cmd/gophkeep/client"
	serverCmd "github.com/Pklerik/gophKeep/cmd/gophkeep/server"
)

// Version and build date (set at build time with -ldflags)
var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "server":
		flag.CommandLine.Parse(os.Args[2:])
		serverCmd.Run()
	case "client":
		flag.CommandLine.Parse(os.Args[2:])
		clientCmd.Run()
	case "version":
		fmt.Printf("GophKeeper version %s (built at %s)\n", Version, BuildTime)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage: gophkeeper <command> [options]

Commands:
  server       Start the GophKeeper server
  client       Start the GophKeeper client
  version      Display version information
  help         Show this help message

Server options:
  -a string    Server address (default "localhost:8080")
  -d string    Database path (default "./gophkeeper.db")
  -s           Enable TLS/HTTPS
  -c string    Certificate file (for TLS)
  -p string    Key file (for TLS)

Client options:
  -u string    Server URL (default "http://localhost:8080")
  -d string    Client database path

Examples:
  gophkeeper server -a 0.0.0.0:8080
  gophkeeper client -u http://localhost:8080
  gophkeeper version`)
}
