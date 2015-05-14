package goatpress

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

type Client struct {
	name   string
	conn   net.Conn
	reader *bufio.Reader
}

func newClient(name string, serverPort int) *Client {
	address := serverAddress + fmt.Sprintf(":%d", serverPort)
	fmt.Printf("Server:      %s\n", address)
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	return &Client{name, conn, bufio.NewReader(conn)}
}

func (c *Client) readLine() string {
	if b, err := c.reader.ReadBytes('\n'); err != nil {
		println(err)
		println("server went away on read")
		os.Exit(0)
	} else {
		line := string(b[0 : len(b)-1])
		fmt.Printf("< %s\n", line)
		return line
	}
	return ""
}

func (c *Client) writeLine(req string) {
	_, err := c.conn.Write([]byte(req + "\n"))
	if err != nil {
		println(err)
		println("server went away on write")
		os.Exit(0)
	}
	fmt.Printf("> %s\n", req)
}

var msgName = regexp.MustCompile(`;\s+name ?`)
var msgNewGame = regexp.MustCompile(`new game`)
var msgPlay = regexp.MustCompile(`; move ([a-z ])+ ([0-9 ])+ ?`)
var msgPing = regexp.MustCompile(`;\s+ping`)

const size = 5

func (c *Client) ParseGetMove(req string) Board {
	bits := strings.Split(req, ";")
	q := bits[1]
	qbits := strings.Split(q, " ")
	boardString := strings.Join(qbits[2:7], "")
	board := newEmptyBoard(size)
	for i, row := range board.Letters {
		for j, _ := range row {
			board.Letters[i][j] = string(boardString[i*size+j])
		}
	}
	return board
}

func (c *Client) Run() {
	moves := make([]Move, 0)
	c.readLine()
	for {
		req := c.readLine()
		if msgNewGame.MatchString(req) {
			moves = make([]Move, 0)
		} else if msgName.MatchString(req) {
			c.writeLine(c.name)
		} else if msgPing.MatchString(req) {
			c.writeLine("pong")
		} else if msgPlay.MatchString(req) {
			board := c.ParseGetMove(req)
			moveFinder := newRandomFinder(DefaultWordSet)
			colorMask := newColorMask(&board, moves)
			state := GameState{0, 0, 0, board, colorMask, colorMask.ToString(), moves}
			move := moveFinder.GetMove(state)
			c.writeLine(move.ToMessage())
		}
	}
}
