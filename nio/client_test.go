// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package nio;

import (
  "testing"
  "log"
  "net/url"
  "github.com/mrichconf/feudal/message"
  "github.com/mrichconf/feudal/worker"
)

func TestNewTcpClient(t *testing.T) {
  u,_ := url.Parse("http://127.0.0.1:5887")
  
  s := &Client{}
  s.Init(s)
  s.Start(worker.DefaultDispatcher(), s)
  log.Println("Sending ConnectTo")
  i := worker.NewInterrogator()
  s.Send(message.New(&ConnectTo{ Url: u }, i))
  r := <- i

  switch r.Body().(type) {
  case *Connected:
  default:
    t.Errorf("no valid response")
  }
  s.Stop()
}


