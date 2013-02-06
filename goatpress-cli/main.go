package main

import (
  "fmt"
  "goatpress"
  "os"
)

func usage() {
  fmt.Printf("usage: goatpress-cli (demo|server|client|web)\n")
}

func main() {
  args := os.Args
  if len(args) < 2 {
    usage()
  } else {
    command := args[1]
    if command == "demo" {
      goatpress.Demo()
    } else if command == "server" {
      goatpress.ServerStart()
    } else if command == "client" {
      goatpress.ClientStart()
    } else if command == "web" {
    } else {
      usage()
    }
  }
}

