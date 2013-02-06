package goatpress

import (
  "net"
  "os"
)

type Client struct {
  Conn net.Conn
}

func newClient() *Client {
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

  return &Client{conn}
}

func (c *Client) Run() {
  _, err := c.Conn.Write([]byte("hello"))
  if err != nil {
    println("Write to server failed:", err.Error())
    os.Exit(1)
  }

  reply := make([]byte, 1024)

  _, err = c.Conn.Read(reply)
  if err != nil {
    println("Read from server failed:", err.Error())
    os.Exit(1)
  }

  println("reply from server=", string(reply))

  c.Conn.Close()
}
