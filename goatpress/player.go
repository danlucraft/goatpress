package goatpress

import (
  "net"
  "fmt"
  "errors"
  "bufio"
  "time"
)

type Player interface {
  Name() string
  NewGame(GameState)
  GetMove(int, string, GameState) Move
  Ping() bool
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

func (p InternalPlayer) Ping() bool {
  return true
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
  req := "; new game \n"
  p.conn.Write([]byte(req))
}

func (p *ClientPlayer) Name() string {
  return p.name
}

func (p *ClientPlayer) Ping() bool {
  p.writeLine("; ping ?")
  line, _ := p.readLine()
  if line != "pong" {
    go p.Unregister()
    p.conn.Close()
    return false
  }
  return true
}

func (p *ClientPlayer) writeLine(req string) error {
  //p.conn.SetWriteDeadline(oneSecondAway())
  _, err := p.conn.Write([]byte(req + "\n"))
  if err != nil {
    go p.Unregister()
    return errors.New("client closed connection on write")
  }
  fmt.Printf("%s> %s\n", p.Name(), req)
  return nil
}

func oneSecondAway() time.Time {
  t := time.Now()
  d, _ := time.ParseDuration("1s")
  t.Add(d)
  return t
}

func (p ClientPlayer) readLine() (string, error) {
  //p.conn.SetReadDeadline(oneSecondAway())
  b, err := p.reader.ReadBytes('\n')
  if err != nil {
    println("client closed connection on read")
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


