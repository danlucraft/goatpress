package goatpress

import (
  "net"
  "os"
  "regexp"
)

type Client struct {
  name string
  conn net.Conn
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

  return &Client{name, conn}
}

func (c *Client) readString() (string,error) {
  buf := make([]byte, 1024)
  if n, err := c.conn.Read(buf); err != nil {
    return "", err
  } else {
    return string(buf[0:n]), nil
  }
  return "", nil
}

var msgNewGame = regexp.MustCompile(`new game`)
var msgPlay    = regexp.MustCompile(`; move ([a-z ]\d)+ ?`)

func (c *Client) Run() {
  version, err := c.readString()
  if err != nil {
    println("couldn't read initial version string")
    os.Exit(1)
  }
  println("version: ", version)
  c.conn.Write([]byte(c.name + "\n"))
  for {
    req, err := c.readString()
    if err != nil {
      println("server went away")
      os.Exit(0)
    }
    println("server said: ", req)
    if msgNewGame.MatchString(req) {
      println("received new game notification")
    } else if msgPlay.MatchString(req) {
      println("play request: sending pass move")
      c.conn.Write([]byte("pass\n"))
    }
  }
}


