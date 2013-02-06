package goatpress

import (
  "net"
  "os"
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

func (c *Client) Run() {
  reply := make([]byte, 1024)
  _, err := c.conn.Read(reply)
  if err != nil {
    println("Read from server failed:", err.Error())
    os.Exit(1)
  }

  println("message from server=", string(reply))

  _, err = c.conn.Write([]byte(c.name))
  if err != nil {
    println("Write to server failed:", err.Error())
    os.Exit(1)
  }

  _, err = c.conn.Read(reply)
}


