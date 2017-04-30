// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package fs

import (
  "os"
)

type IFileSystem interface {
  Open(name string) (IFile, error)
  Stat(name string) (os.FileInfo, error)
}

// Local filesystem access

type localFileSystem struct{}

func (localFileSystem) Open(name string) (IFile, error) { return os.Open(name) }
func (localFileSystem) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

