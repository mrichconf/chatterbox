// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package pbserv;

import (
  "io"
  "github.com/mrichconf/feudal/message"
  "github.com/mrichconf/feudal/worker"
  "github.com/mrichconf/chatterbox/nio"
)

// Handles Messages associated with an individual connection
type ConnectionWorker struct {
  AbstractWorker
  stream worker.Runner
  
  handler worker.Factory
  target worker.Context
}

func (cw *ConnectionWorker) Receive(m message.Envelope) {
  switch r := m.Body().(type) {
  case message.Register:
    if r.handler == nil {
      cw.target = r.Sender()
    } else {
      cw.handler = r.handler
      cw.target = cw.WorkerWhence(func() worker.Context { return cw.handler() })
    }
  case *nio.Connection:
    cw.stream = NewProtocol(r.Conn)
    cw.stream.Start(worker.DefaultDispatcher(), cw) // start the read loop
  case *message.Serialized:
    cw.target.Send(m)
  case *nio.WriteRequest: // generally the registered handler would send us WriteRequest messages
    // write to socket
    n, e := cw.stream.Conn.Write(b.Content)
    s.Send(message.New(&WriteResponse{ Size: n, Origin: b.Origin }, c))
  }
  
}
