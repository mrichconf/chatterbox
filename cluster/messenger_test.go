// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package cluster

import (
  "testing"
  "github.com/hashicorp/memberlist"
  "github.com/golang/protobuf/proto"
  "github.com/golang/protobuf/proto/testdata"
  "fmt"
  "time"
)

func TestNewMessenger(t *testing.T) {
  m := NewMessenger(func() int { return 1 }, 2)
  if m == nil {
    t.Errorf("Failed creating new messenger")
  }
}

func TestBroadcast(t *testing.T) {
  config4 := memberlist.DefaultLocalConfig()
  config4.Name = "config4"
  config4.BindPort = 7950

  c4 := New(config4)


  host := []string{"127.0.0.1:7947" }
  c4.Join(host)

  for _, member := range c4.list.Members() {
    fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
  }
  m := c4.GetMessenger()

  config5 := memberlist.DefaultLocalConfig()
  config5.Name = "config5"
  config5.BindPort = 7951

  c5 := New(config5)

  host = []string{"127.0.0.1:7950" }
  c5.Join(host)

  test := &testdata.GoTestField{ Label: proto.String("test"), Type: proto.String("type") }
  e := m.Queue(test)
  if m.Len() < 1 || e != nil {
    t.Errorf("Failed to queue message", e)
  }

  m5 := c5.GetMessenger()
  res := &testdata.GoTestField{}
  er := m5.Receive(res)

  if er != nil || res == nil {
    t.Errorf("Failed receiving broadcast message: ", er)
  }
  if *res.Label != "test" {
    t.Errorf("Received message doesn't match sent message")
  }
  time.Sleep(1 * time.Second)

  fmt.Println("Received proto message: ", res)
}
