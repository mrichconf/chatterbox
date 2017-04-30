// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package nio;

import (
  "log"
  "net"
  "github.com/mrichconf/feudal/message"
  "github.com/mrichconf/feudal/worker"
)

// Implements Runner interface
type Socket struct {
  listener net.Listener
  quit chan bool
  stopped chan bool
}

func NewSocket(l net.Listener) *Socket {
  return &Socket{ listener: l, quit: make(chan bool), stopped: make(chan bool) }
}

func (s *Socket) Start(d worker.Dispatcher, c worker.Context) {
  go func() {
    for {
      select {
      case <- s.quit:
        s.listener.Close()
        //q.Dispatch(message.New(&message.Terminated{}, nil), c)
        if v := c.Division(); v != nil {
          v.Disolve()
        }
        close(s.stopped)
        return
      default:
        x, e := s.listener.Accept()
        if e != nil {
          log.Fatal("Error accepting: ", e.Error())
        }
        //w := c.WorkerWhence(func() worker.Context { return c.handler() })
        c.Send(message.New(&Connection{ Conn: x }, c))
      }
    }
  }()
}

func (s *Socket) Stop() {
  close(s.quit)
  select {
  case <- s.stopped:
  }
  log.Println("Terminating socket runner")
}  
