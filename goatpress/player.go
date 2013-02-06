package goatpress

import (
  "net"
)

type Player interface {
  Name() string
  NewGame(GameState)
  GetMove(int, string, GameState) Move
}

const (
  MSG_NONE               = iota
  MSG_BAD_MOVE_UNKNOWN   = iota
  MSG_BAD_MOVE_PREFIX    = iota
  MSG_BAD_MOVE_ALREADY   = iota
  MSG_BAD_MOVE_TOO_SHORT = iota
  MSG_OPPONENT_MOVE      = iota
)

type InternalPlayer struct {
  name       string
  MoveFinder MoveFinder
}

func newInternalPlayer(name string, moveFinder MoveFinder) InternalPlayer {
  return InternalPlayer{name, moveFinder}
}

func (p InternalPlayer) NewGame(_ GameState) {
}

func (p InternalPlayer) Name() string {
  return p.name
}

func (p InternalPlayer) GetMove(_ int, _ string, state GameState) Move {
  moveFinder := p.MoveFinder
  return moveFinder.GetMove(state)
}

type ClientPlayer struct {
  name       string
  conn       net.Conn
  unregister chan string
}

type ClientMessage struct {
  msgType int
  info    string
  request string
}

func (p ClientPlayer) NewGame(_ GameState) {
  req := "new game; \n"
  p.conn.Write([]byte(req))
}

func (p ClientPlayer) Name() string {
  return p.name
}

func (p ClientPlayer) readString() string {
  buf := make([]byte, 1024)
  n, _ := p.conn.Read(buf)
  return string(buf[0:n])
}

func newClientPlayer(conn net.Conn, unregister chan string) ClientPlayer {
  conn.Write([]byte("name?\n"))
  buf := make([]byte, 1024)
  n, _ := conn.Read(buf)
  name := string(buf[0:n])
  return ClientPlayer{name, conn, unregister}
}

func (p ClientPlayer) GetMove(msg int, info string, state GameState) Move {
  board := state.Board.ToString()
  var bit1 string
  switch msg {
    case MSG_NONE:               bit1 = ""
    case MSG_BAD_MOVE_UNKNOWN:   bit1 = "invalid: unknown"
    case MSG_BAD_MOVE_PREFIX:    bit1 = "invalid: prefix"
    case MSG_BAD_MOVE_TOO_SHORT: bit1 = "invalid: too-short"
    case MSG_OPPONENT_MOVE:      bit1 = "opponent: " + info
  }
  req := bit1 + " ; play '" + board + "'\n"
  p.conn.Write([]byte(req))
  data := p.readString()
  println(data)
  return MakePassMove()
}


