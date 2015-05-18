package goatpress

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Player interface {
	Name() string
	NewGame(GameState, string)
	GetMove(int, string, GameState) Move
	Disconnect(string)
	Ping() bool
}

const (
	MSG_NONE = iota
	MSG_BAD_MOVE_UNKNOWN
	MSG_BAD_MOVE_PREFIX
	MSG_BAD_MOVE_ALREADY
	MSG_BAD_MOVE_TOO_SHORT
	MSG_OPPONENT_MOVE
)

type InternalPlayer struct {
	name       string
	MoveFinder MoveFinder
}

func newInternalPlayer(name string, moveFinder MoveFinder) InternalPlayer {
	return InternalPlayer{name, moveFinder}
}

func (p InternalPlayer) NewGame(_ GameState, OpponentName string) {
}

func (p InternalPlayer) Name() string {
	return p.name
}

func (p InternalPlayer) Disconnect(_ string) {
	return
}

func (p InternalPlayer) Ping() bool {
	return true
}

func (p InternalPlayer) GetMove(_ int, _ string, state GameState) Move {
	moveFinder := p.MoveFinder
	return moveFinder.GetMove(state)
}

type ClientPlayer struct {
	name          string
	conn          net.Conn
	unregister    chan string
	reader        *bufio.Reader
	clientTimeout string
	closed        bool
}

type ClientMessage struct {
	msgType int
	info    string
	request string
}

func newClientPlayer(conn net.Conn, unregister chan string, clientTimeout string) *ClientPlayer {
	p := &ClientPlayer{"", conn, unregister, bufio.NewReader(conn), clientTimeout, false}
	for p.name == "" {
		err := p.writeLine("; name ?")
		if err != nil {
			return nil
		}
		str, err2 := p.readLine()
		if err2 != nil {
			return nil
		}
		if ValidateName(str) {
			p.name = str
		} else {
			p.writeLine("invalid name")
			conn.Close()
			return nil
		}
	}
	return p
}

func (p *ClientPlayer) NewGame(_ GameState, OpponentName string) {
	p.writeLine(fmt.Sprintf("new game vs '%s';", OpponentName))
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
	if p.closed {
		return errors.New("closed")
	}
	p.conn.SetWriteDeadline(p.deadline())
	_, err := p.conn.Write([]byte(req + "\n"))
	if err != nil {
		fmt.Println(err)
		go p.Unregister()
		return errors.New("client closed connection on write")
	}
	//fmt.Printf("%s> %s\n", p.Name(), req)
	return nil
}

var nameValidate = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`)

func ValidateName(name string) bool {
	return nameValidate.MatchString(name)
}

func (p ClientPlayer) deadline() time.Time {
	t := time.Now()
	d, _ := time.ParseDuration("10s")
	return t.Add(d)
}

func (p ClientPlayer) readLine() (string, error) {
	if p.closed {
		return "", errors.New("closed")
	}
	p.conn.SetReadDeadline(p.deadline())
	b, err := p.reader.ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
		go p.Unregister()
		return "", errors.New("client closed connection")
	}

	var lastByte byte
	if len(b) > 0 {
		lastByte := b[len(b)-1]
		if lastByte == 10 || lastByte == 13 {
			b = b[0 : len(b)-1]
		}
	}

	if len(b) > 0 {
		lastByte = b[len(b)-1]
		if lastByte == 10 || lastByte == 13 {
			b = b[0 : len(b)-1]
		}
	}

	line := string(b[0:len(b)])
	//fmt.Printf("%s< %s\n", p.Name(), line)
	return line, nil
}

func (p *ClientPlayer) Disconnect(reason string) {
	p.writeLine(fmt.Sprintf("unregistering: %s", reason))
	p.conn.Close()
	p.closed = true
}

func (p *ClientPlayer) Unregister() {
	p.Disconnect("\"reasons\"")
	p.unregister <- p.name
}

var moveFormat = regexp.MustCompile(`^move:([0-9][0-9],?)+ \([^\)]*\)$`)

func (p *ClientPlayer) GetMove(msg int, info string, state GameState) Move {
	board := state.Board.ToString()
	colors := state.ColorString
	var bit1 string
	switch msg {
	case MSG_NONE:
		bit1 = ""
	case MSG_BAD_MOVE_ALREADY:
		bit1 = "invalid: already"
	case MSG_BAD_MOVE_UNKNOWN:
		bit1 = "invalid: unknown"
	case MSG_BAD_MOVE_PREFIX:
		bit1 = "invalid: prefix"
	case MSG_BAD_MOVE_TOO_SHORT:
		bit1 = "invalid: too-short"
	case MSG_OPPONENT_MOVE:
		bit1 = "opponent: " + info
	}
	req := bit1 + " ; move " + board + " " + colors + " ?"
	err1 := p.writeLine(req)
	if err1 != nil {
		//fmt.Printf("%s passes due to closed connection\n", p.name)
		return MakePassMove()
	}

	data, err2 := p.readLine()
	if err2 != nil {
		//fmt.Printf("%s passes due to closed connection\n", p.name)
		return MakePassMove()
	}

	if data == "pass" {
		return MakePassMove()
	} else if moveFormat.MatchString(data) {
		bits := strings.Split(data, ":")
		moveString := bits[1]
		tileStrings := strings.Split(moveString, ",")
		tiles := make([]Tile, len(tileStrings))
		for i, ts := range tileStrings {
			xi, _ := strconv.ParseInt(string(ts[0]), 10, 0)
			yi, _ := strconv.ParseInt(string(ts[1]), 10, 0)
			tile := newTile(int(xi), int(yi))
			tiles[i] = tile
		}
		move := state.Board.MoveFromTiles(tiles)
		return move
	} else {
		println("bad data '", data, "'")
		p.writeLine("invalid: bad-format, passing ; ")
		return MakePassMove()
	}

	return MakePassMove()
}

func dummyFtm() {
	fmt.Printf("")
}
