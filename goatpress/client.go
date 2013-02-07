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

func (c *Client) readLine() (string,error) {
  if b, err := c.reader.ReadBytes('\n'); err != nil {
    return "", err
  } else {
    return string(b[0:len(b)-1]), nil
  }
  return "", nil
}

var msgNewGame = regexp.MustCompile(`new game`)
var msgPlay    = regexp.MustCompile(`; move ([a-z 0-9])+ ?`)

func (c *Client) Run() {
  version, err := c.readLine()
  if err != nil {
    println("couldn't read initial version string")
    os.Exit(1)
  }
  println("version: ", version)
  c.conn.Write([]byte(c.name + "\n"))
  for {
    req, err := c.readLine()
    if err != nil {
      println("server went away")
      os.Exit(0)
    }
    println("server said: ", req)
    if msgNewGame.MatchString(req) {
      println("received new game notification")
    } else if msgPlay.MatchString(req) {
      println("play request: sending pass move")
      n, err := c.conn.Write([]byte("pass\n"))
      fmt.Printf("wrote %d bytes\n", n)
      if err != nil {
        fmt.Printf("error writing pass\n")
      }
    }
  }
}


