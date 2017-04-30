// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package nio;

import (
  "io"
  "github.com/mrichconf/feudal/worker"
)

type Streamer interface {
  io.Closer
  io.Reader
  io.Writer
}

// messages

type Connection struct {
  Conn Streamer
}

type ReadRequest struct { // WriteReponse
  Size uint32
  Origin worker.Ref
}
type WriteResponse ReadRequest

type ReadResponse struct { // WriteRequest
  Content []byte
  Origin worker.Ref
}
type WriteRequest ReadResponse
