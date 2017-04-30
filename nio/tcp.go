// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package nio;

import (
  "net"
  "net/url"
  "log"
  "github.com/mrichconf/feudal/worker"
  "github.com/mrichconf/feudal/message"
)

// TCP Listener
// create workers for each connection and manage them

type Listen struct {
  Url *url.URL
}
type Accepting struct {}

type Tcp struct {
  worker.AbstractContext
  socket worker.Runner
  addr *url.URL
  handler worker.Factory
}

func (t *Tcp) Receive(m message.Envelope) {
  s := m.Sender()
  switch b := m.Body().(type) {
//  case *worker.Terminated:
  case *message.Register:
    if b.Handler == nil {
      //t.target = b.Sender()
      //error
    } else {
      t.handler = b.Handler.(worker.Factory)
    }
  case *Listen:
    l, e := net.Listen("tcp", b.Url.Host)
    if e != nil {
      log.Fatal("Error listening: ", e.Error())
    }
    log.Println("Opening socket for listen request for: ", b.Url.Host)
    t.socket = NewSocket(l)
    t.socket.Start(worker.DefaultDispatcher(), t)
    if s != nil {
      s.Send(message.New(&Accepting{}, t))
    }
  case *Connection:
    w := t.WorkerWhence(func() worker.Context { return t.handler() })
    w.Send(message.New(b, t))
  default:
  }
}

func (t *Tcp) Stop() {
  t.socket.Stop()
  t.AbstractContext.Queue.Stop()
}

/*
func (s *Socket) Start(c worker.Context) {
  go func() {
    defer s.listener.Close()
//    pbs.connection.(*net.TCPListener).SetDeadline(time.Now().Add(time.Second))
    for {
      select {
      case r := <- s.Queue.quit:
          log.Println("Listen(): quitting")
          
          return
      default:
        x, e := s.listener.Accept()
        if e != nil {
          log.Fatal("Error accepting: ", e.Error())
        }
        w := c.WorkerWhence(func() worker.Context { return c.handler() })
        w.Send(message.New(&Connection{ Conn: x }, c))
      }
    }
  }()
}

func (s *Socket) Stop() {
  s.(*Queue).Stop()
}
*/
