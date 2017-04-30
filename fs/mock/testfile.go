// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package mock_testfile;

import (
  "github.com/mrichconf/chatterbox/fs"
  "io"
  "os"
  "errors"
  //"log"
)

type TestFileSystem struct{}

func (TestFileSystem) Open(name string) (fs.IFile, error) { return os.Open(name) }
func (TestFileSystem) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

type TestFile struct {
  data *[]byte
  offset *int64
}

func New() *TestFile {
  return &TestFile{ data: new([]byte), offset: new(int64) }
}

func (t TestFile) Write(p []byte) (n int, err error) {
  *t.data = append(*t.data, p...)
  *t.offset += int64(len(p))
  if int64(len(*t.data)) != *t.offset {
    return len(p), errors.New("write error")
  }
  return len(p), nil
}

func (t TestFile) Read(p []byte) (n int, err error) {
  if *t.offset >= int64(len(*t.data)) {
    return 0, io.EOF
  }
  n = copy(p, (*t.data)[*t.offset:])
  *t.offset += int64(n)
  if n < len(p) {
    return n, io.EOF
  }
  return n, nil
}

func (t TestFile) ReadAt(p []byte, off int64) (n int, err error) {
  *t.offset = off
  return t.Read(p)
}

func (t TestFile) Close() error {
  *t.data = nil
  *t.offset = 0
  return nil
}

func (t TestFile) Seek(offset int64, whence int) (int64, error) {
  var o int64 = 0
  switch whence {
  case 0:
    o = offset
  case 1:
    fallthrough
  case 2:
    o = *t.offset + offset
  }

  if o > int64(len(*t.data)) || o < 0 {
    return 0, &os.PathError{"seek", "", errors.New("invalid offset")}
  }
  *t.offset = o 
  return *t.offset, nil
}

func (t TestFile) Stat() (os.FileInfo, error) {
  var info os.FileInfo
  return info, nil
}
