// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package block;

import (
  "github.com/mrichconf/feudal/message"
  "github.com/mrichconf/feudal/worker"
)

// communicates with a block storage service

type ClientWorker struct {
  worker.AbstractWorker
  protobuf worker.Ref // pb.Codec <-> nio.Client
  connection worker.Ref
}

func (c *ClientWorker) Receive(m message.Envelope) {
  s := m.Sender()
  switch b := m.Body().(type) {
  case *ConnectTo:
    c.connection = c.WorkerWhence(func() worker.Context { return &nio.Client{} })
    c.connection.Send(m)
    c.protobuf = c.WorkerWhence(func() worker.Context { return &pb.Codec{ client: c.connection, messageFactory: func() worker.Context { return &Block{} } } })
  case *BlockAssertion:
    // send block
    c.protobuf.Send(message.New(&Encode{ Message: b.Block }, s)) // response will go to our sender
  }
}

