// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package fs

import (
  "io"
  "os"
)

type IFile interface {
  io.Closer
  io.Reader
  io.ReaderAt
  io.Seeker
  io.Writer
  Stat() (os.FileInfo, error)
  Truncate(size int64) error
}
