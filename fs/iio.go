// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package fs

import (
  "io"
//  "os"
)

type IIO interface {
  io.Closer
  io.Reader
  io.Writer
}
