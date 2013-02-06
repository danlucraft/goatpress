package goatpress

import (
  "net"
)

type Player interface {
  Name() string
  GetMove(GameState) Move
}

type InternalPlayer struct {
  name       string
  MoveFinder MoveFinder
}

func newInternalPlayer(name string, moveFinder MoveFinder) InternalPlayer {
  return InternalPlayer{name, moveFinder}
}

func (p InternalPlayer) Name() string {
  return p.name
}

func (p InternalPlayer) GetMove(state GameState) Move {
  moveFinder := p.MoveFinder
  return moveFinder.GetMove(state)
}

type ClientPlayer struct {
  name string
  conn net.Conn
}

func (p ClientPlayer) Name() string {
  return p.name
}

func newClientPlayer(conn net.Conn) ClientPlayer {
  buf := make([]byte, 1024)
  n, _ := conn.Read(buf)
  name := string(buf[0:n])
  return ClientPlayer{name,conn}
}

func (p ClientPlayer) GetMove(state GameState) Move {
  return MakePassMove()

}

