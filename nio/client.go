// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package nio;

import (
  "net"
  "io"
  "net/url"
  "log"
  "github.com/mrichconf/feudal/worker"
  "github.com/mrichconf/feudal/message"
)

// TCP Listener
// create workers for each connection and manage them

type ConnectTo struct {
  Url *url.URL
}
type Connected struct {}

// will need port to listen on and worker factory to handle connections
type Client struct {
  worker.AbstractContext
  addr *url.URL
  connection *Connection
}

func (c *Client) read(s int, b []byte) error {
  const errorLimit = 3
  readSlice := b[:]
  a := 0
  z := s
  ioErrors uint32 = 0
  while {
    n,e := io.ReadAtLeast(c.connection.Conn, readSlice, z)
    if n == z {
      return nil
    }
    if n > 0 && n < z {
      // io.ErrUnexpectedEOF or something else
      a = n
      z -= n
      readSlice = readSlice[a:z]
    }
    if e == io.EOF { // closed
      return nil
    }
    if e != nil {
      if ++ioErrors > errorLimit {
        // unrecoverable - we can return this error and it will get passed back to the manager
        return e
      }
    }
  }
}

// tcp client recieves messages and proxies them to the tcp connection
// need a runner to read/proccess responses?
func (c *Client) Receive(m message.Envelope) {
  s := m.Sender()
  switch b := m.Body().(type) {
  //case *message.Terminated:
  case *ConnectTo:
    c.addr = b.Url
    n, e := net.Dial("tcp", c.addr.Host)
    if e != nil {
      log.Fatal("Error listening: ", e.Error())
    }
    c.connection = &Connection{ Conn: n }
    if s != nil {
      s.Send(message.New(&Connected{}, c))
    }
  case *WriteRequest:
    n, e := c.connection.Conn.Write(b.Content)
    s.Send(message.New(&WriteResponse{ Size: uint32(n), Origin: b.Origin }, c))
  case *ReadRequest: // get response TODO: fix the naming
    buffer := make([]byte, b.Size)
    if e := c.read(b.Size, buffer); e == nil {
      s.Send(message.New(&ReadResponse{ Content: buffer, Origin: b.Origin }, c))
    } else {
      c.Manager().Send(message.New(&message.Error{ e: e }), c)
      // or should this be
      //c.manager.TerminateWorker()
    }
  }
}

func (c *Client) Stop() {
  c.AbstractContext.Queue.Stop()
  c.connection.Close()
}

