// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package block;

import (
  "github.com/mrichconf/feudal/message"
  "github.com/mrichconf/feudal/worker"
)

// An assertion worker resolves an assertion to a stored block (if it exists) or the assertion fails

type Cache interface {
  Get(ctx groupcache.Context, key string, dest groupcache.Sink) error
  CacheStats(which groupcache.CacheType) groupcache.CacheStats
}

type AssertionWorker struct {
  worker.AbstractWorker
  blocks Cache
  blockStorageClient worker.Ref 
}

func (bw *AssertionWorker) Receive(m message.Envelope) {
  switch b := m.Body().(type) {
  case *BlockAssertion:

    // request from cache
    blocks.Get(nil, b.GetId(), d)
   
    bw.blockStorageClient.Send(message.New(b, bw.Sender())) // the storage client should respond directly to the sender

   // send assertion response
  }
}
