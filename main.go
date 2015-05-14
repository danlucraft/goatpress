package main

import (
	"flag"
	"fmt"
	"github.com/danlucraft/goatpress/goatpress"
)

func usage() {
	fmt.Printf("usage: goatpress-cli OPTIONS (demo|server|client)\n")
	fmt.Printf("\n")
	fmt.Printf("  -name=NAME          Player Name (if running as client)\n")
	fmt.Printf("  -file=PATH          Where to store tournament history (if server)\n")
	fmt.Printf("  -serverport=PORT    Game server port (if server/client) default 4123\n")
	fmt.Printf("  -webport=PORT       Web page port (if server) default 5123\n")
}

func main() {
	playerName := flag.String("name", "", "name")
	tournamentFile := flag.String("file", "", "file")
	clientTimeout := flag.String("timeout", "1s", "timeout")
	serverPort := flag.Int("serverport", 4123, "serverport")
	webPort := flag.Int("webport", 5123, "webport")

	flag.Parse()

	if len(flag.Args()) < 1 {
		usage()
	} else {
		fmt.Printf("Running as:  %s\n", flag.Args()[0])
		command := flag.Args()[0]
		if command == "demo" {
			goatpress.Demo()
		} else if command == "server" {
			fmt.Printf("Web port:    %s\n", *webPort)
			goatpress.ServerStart(*tournamentFile, *clientTimeout, *serverPort, *webPort)
		} else if command == "client" {
			fmt.Printf("Player name: %s\n", *playerName)
			goatpress.ClientStart(*playerName, *serverPort)
		} else {
			usage()
		}
	}
}
