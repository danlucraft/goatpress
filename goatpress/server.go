package goatpress

import (
  "net"
  "os"
  "fmt"
  "time"
)

const serverAddress = "localhost:4123"
var newPlayers = make(chan Player)

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
  go AcceptPlayers(listener)
  for {
    fmt.Printf("select:\n")
    select {
    case newPlayer:= <- newPlayers:
      fmt.Printf("Player Online: %s\n", newPlayer.Name())
      c.Tournament.RegisterPlayer(newPlayer)
    default:
      if c.Tournament.Size() > 1 {
        c.Tournament.PlayMatch()
      }
    }
    time.Sleep(1e9)
  }
}

const serverSig = "goatpress<VERSION=1>;"

func AcceptPlayers(listener net.Listener) {
  for {
    println("waiting")
    conn, err := listener.Accept()
    if err != nil {
      println("Error accept:", err.Error())
      return
    }
    conn.Write([]byte(serverSig))
    player := newClientPlayer(conn)
    newPlayers <- player
  }
}


