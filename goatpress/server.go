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

func EchoFunc(conn net.Conn) {
  println("Connection!")
  buf := make([]byte, 1024)
  conn.Read(buf)
  println(string(buf))
  conn.Write([]byte("bakc"))
}

func (c *Server) Run() {
  listener, err := net.Listen("tcp", serverAddress)
  if err != nil {
    fmt.Printf("error listening:", err.Error())
    os.Exit(1)
  }

  for {
    println("listening")
    conn, err := listener.Accept()
    if err != nil {
      println("Error accept:", err.Error())
      return
    }
    go EchoFunc(conn)
  }
}
