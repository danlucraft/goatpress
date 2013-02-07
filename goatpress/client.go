package goatpress

import (
  "net"
  "os"
  "regexp"
  "bufio"
  "fmt"
)

type Client struct {
  name   string
  conn   net.Conn
  reader *bufio.Reader
}

func newClient(name string) *Client {
  servAddr := "localhost:4123"
  tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
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
    line := string(b[0:len(b)-1])
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

var msgName    = regexp.MustCompile(`;\s+name ?`)
var msgNewGame = regexp.MustCompile(`; new game`)
var msgPlay    = regexp.MustCompile(`; move ([a-z 0-9])+ ?`)
var msgPing    = regexp.MustCompile(`;\s+ping`)

func (c *Client) Run() {
  c.readLine()
  for {
    req := c.readLine()
    if msgNewGame.MatchString(req) {
    } else if msgName.MatchString(req) {
      c.writeLine(c.name)
    } else if msgPing.MatchString(req) {
      c.writeLine("pong")
    } else if msgPlay.MatchString(req) {
      c.writeLine("pass")
    }
  }
}




