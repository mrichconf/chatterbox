// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package cluster

import (
  "testing"
  "github.com/hashicorp/memberlist"
  "fmt"
  "time"
)

func TestNewCluster(t *testing.T) {
  config := memberlist.DefaultLocalConfig()
  config.BindPort = 7946
  c := New(config)
  if c == nil {
    t.Errorf("Failed creating new cluster")
  }
}

func TestJoinCluster(t *testing.T) {
  config1 := memberlist.DefaultLocalConfig()
  config2 := memberlist.DefaultLocalConfig()

  config1.Name = "config1"
  config1.BindPort = 7947
  config2.Name = "config2"
  config2.BindPort = 7948

  c1 := New(config1)
  if c1 == nil {
    t.Errorf("Failed creating new cluster")
  }

  c2 := New(config2)
  if c2 == nil {
    t.Errorf("Failed creating new cluster")
  }

  host := []string{"127.0.0.1:7947", "127.0.0.1:7946" }
  c2.Join(host)

  fmt.Println("successfully contacted: ")

  for _, member := range c2.list.Members() {
    fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
  }
  time.Sleep(1 * time.Second)
}

func TestClusterEvents(t *testing.T) {
  joinCh := make(chan *memberlist.Node, 50)
  config3 := memberlist.DefaultLocalConfig()
  config3.Name = "config3"
  config3.BindPort = 7949

  c3 := New(config3)
  
  ce := c3.GetEvents()
  ce.RegisterJoinChannel(joinCh)

  host := []string{"127.0.0.1:7947" }
  c3.Join(host)

  res := <- joinCh
  fmt.Println(res)

  for _, member := range c3.list.Members() {
    fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
  }
  fmt.Println("joinCh: ", len(joinCh))
}
