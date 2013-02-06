package main

import (
  "fmt"
  "goatpress"
  "os"
)

func main() {
  args := os.Args
  if len(args) < 2 {
    fmt.Printf("usage: goatpress-cli (DEMO)\n")
  } else {
    command := args[1]
    if command == "demo" {
      goatpress.Demo()
    }
  }
}

