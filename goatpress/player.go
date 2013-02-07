package goatpress

import (
  "net"
  "fmt"
  "errors"
  "bufio"
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
  reader     *bufio.Reader
}

type ClientMessage struct {
  msgType int
  info    string
  request string
}

func newClientPlayer(conn net.Conn, unregister chan string) *ClientPlayer {
  p := &ClientPlayer{"", conn, unregister, bufio.NewReader(conn)}
  for p.name == "" {
    err := p.writeLine("; name ?")
    if err != nil {
      return nil
    }
    str, err2 := p.readLine()
    if err2 != nil {
      return nil
    }
    if len(str) > 0 {
      p.name = str
    }
  }
  return p
}

func (p *ClientPlayer) NewGame(_ GameState) {
  req := "new game; \n"
  p.conn.Write([]byte(req))
}

func (p *ClientPlayer) Name() string {
  return p.name
}

func (p *ClientPlayer) writeLine(req string) error {
  _, err := p.conn.Write([]byte(req))
  if err != nil {
    go p.Unregister()
    return errors.New("client closed connection")
  }
  fmt.Printf("%s> %s\n", p.Name(), req[0:len(req)-1])
  return nil
}

func (p ClientPlayer) readLine() (string, error) {
  b, err := p.reader.ReadBytes('\n')
  if err != nil {
    go p.Unregister()
    return "", errors.New("client closed connection")
  }
  line := string(b[0:len(b)-1])
  fmt.Printf("%s< %s\n", p.Name(), line)
  return line, nil
}

func (p *ClientPlayer) Unregister() {
  p.unregister <- p.name
}

func (p *ClientPlayer) GetMove(msg int, info string, state GameState) Move {
  board := state.Board.ToString()
  colors := state.ColorString
  var bit1 string
  switch msg {
    case MSG_NONE:               bit1 = ""
    case MSG_BAD_MOVE_UNKNOWN:   bit1 = "invalid: unknown"
    case MSG_BAD_MOVE_PREFIX:    bit1 = "invalid: prefix"
    case MSG_BAD_MOVE_TOO_SHORT: bit1 = "invalid: too-short"
    case MSG_OPPONENT_MOVE:      bit1 = "opponent: " + info
  }
  req := bit1 + " ; move " + board + " " + colors + " ?\n"
  err1 := p.writeLine(req)
  if err1 != nil {
    fmt.Printf("%s passes due to closed connection\n", p.name)
    return MakePassMove()
  }

  data, err2 := p.readLine()
  if err2 != nil {
    fmt.Printf("%s passes due to closed connection\n", p.name)
    return MakePassMove()
  }
  if data == "pass\n" {
    return MakePassMove()
  }
  return MakePassMove()
}


