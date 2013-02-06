package goatpress

import (
  "net"
  "os"
  "fmt"
)

const serverAddress = "localhost:4123"

type Server struct {
  Tournament *Tournament
}

func newServer() *Server {
  gameType := newGameType(5, DefaultWordSet)
  tourney := newTournament(*gameType)
  return &Server{tourney}
}

func (c *Server) Run() {
  listener, err := net.Listen("tcp", serverAddress)
  if err != nil {
    fmt.Printf("error listening:", err.Error())
    os.Exit(1)
  }

  for {
    println("waiting")
    conn, err := listener.Accept()
    if err != nil {
      println("Error accept:", err.Error())
      return
    }
    go AcceptPlayer(conn)
  }
}

const serverSig = "goatpress<VERSION=1>;"

func AcceptPlayer(conn net.Conn) {
  conn.Write([]byte(serverSig))
  player := newClientPlayer(conn)
  fmt.Printf("Player Online: %s\n", player.Name())
}


