package main

import (
	"flag"
	"fmt"
	"github.com/danlucraft/goatpress/goatpress"
)

func usage() {
	fmt.Printf("usage: goatpress-cli OPTIONS (demo|server|client|web)\n")
}

func main() {
	playerName := flag.String("name", "", "name")
	tournamentFile := flag.String("file", "", "file")
	clientTimeout := flag.String("timeout", "1s", "timeout")
	serverPort := flag.Int("serverport", 4123, "serverport")
	webPort := flag.Int("webport", 5123, "webport")
	flag.Parse()
	fmt.Printf("command %s\n", flag.Args()[0])
	if len(flag.Args()) < 1 {
		usage()
	} else {
		command := flag.Args()[0]
		if command == "demo" {
			goatpress.Demo()
		} else if command == "server" {
			goatpress.ServerStart(*tournamentFile, *clientTimeout, *serverPort, *webPort)
		} else if command == "client" {
			fmt.Printf("Connecting as player %s\n", *playerName)
			goatpress.ClientStart(*playerName, *serverPort)
		} else if command == "web" {
		} else {
			usage()
		}
	}
}
