// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package cluster

import (
  "errors"
  "github.com/golang/protobuf/proto"
  "github.com/hashicorp/memberlist"
//  "github.com/mrichconf/chatterbox/cluster/proto"
  "github.com/mrichconf/feudal/message"
)

type ProtobufBroadcaster interface {
  Queue(msg proto.Message) error
}

type ProtobufReceiver interface {
  Receive(proto.Message) error
}

type MessageInjector interface {
  NewMessage(r ProtobufReceiver) error
  ProtobufReceiver
}

type ProtobufMessenger interface {
  ProtobufBroadcaster
  ProtobufReceiver
  RegisterMessageInjector(i MessageInjector)

  RegisterMessageType(mType uint32, msg message.Factory) error
  DecodeMessage(mType uint32, data []byte) (proto.Message, error)
  Len() int
}

type Cluster struct {
  config *memberlist.Config
  list *memberlist.Memberlist
}

type ClusterEvents struct {
  joinEvents []chan *memberlist.Node
  leaveEvents []chan *memberlist.Node
  updateEvents []chan *memberlist.Node
  passthru memberlist.EventDelegate
}

func New(config *memberlist.Config) *Cluster {
  if config == nil {
    config = memberlist.DefaultLocalConfig()
  }
  l, err := memberlist.Create(config)
  if err != nil {
    panic("Failed to create memberlist: " + err.Error())
  }

  ce := &ClusterEvents{ passthru: config.Events }
  config.Events = ce
  c :=  &Cluster{ config: config, list: l }

  m := NewMessenger(func() int { return c.NumMembers() }, config.RetransmitMult)
  m.Receiver()
  c.config.Delegate = m
  return c
}

func (c *Cluster) GetHost() string {
  return c.config.BindAddr
}

func (c *Cluster) Members() []*memberlist.Node {
  return c.list.Members()
}

func (c *Cluster) NumMembers() int {
  return c.list.NumMembers()
}

func (c *Cluster) Join(ip []string) (int, error) {
  n, err := c.list.Join(ip)
  if err != nil {
    panic("Failed to join cluster: " + err.Error())
  }
  return n, err
}

func (c *Cluster) GetMessenger() ProtobufMessenger {
  return c.config.Delegate.(ProtobufMessenger)
}

func (c *Cluster) GetEvents() *ClusterEvents {
  return c.config.Events.(*ClusterEvents)
}

func (ce *ClusterEvents) RegisterJoinChannel(join chan *memberlist.Node) error {
  if cap(join) == 0 {
    return errors.New("Only a buffered channel can be registerd for join events")
  }
  ce.joinEvents = append(ce.joinEvents, join)
  return nil
}

func (ce *ClusterEvents) RegisterLeaveChannel(leave chan *memberlist.Node) error {
  if cap(leave) == 0 {
    return errors.New("Only a buffered channel can be registerd for leave events")
  }
  ce.leaveEvents = append(ce.leaveEvents, leave)
  return nil
}

func (ce *ClusterEvents) RegisterUpdateChannel(update chan *memberlist.Node) error {
  if cap(update) == 0 {
    return errors.New("Only a buffered channel can be registerd for update events")
  }
  ce.updateEvents = append(ce.updateEvents, update)
  return nil
}

// Event delegate interface
func (ce *ClusterEvents) NotifyJoin(n *memberlist.Node) {
  for _, ch := range ce.joinEvents {
    ch <- n
  }
  if ce.passthru != nil {
    ce.passthru.NotifyJoin(n)
  }
}

func (ce *ClusterEvents) NotifyLeave(n *memberlist.Node) {
  for _, ch := range ce.leaveEvents {
    ch <- n
  }
  if ce.passthru != nil {
    ce.passthru.NotifyLeave(n)
  }
}

func (ce *ClusterEvents) NotifyUpdate(n *memberlist.Node) {
  for _, ch := range ce.updateEvents {
    ch <- n
  }
  if ce.passthru != nil {
    ce.passthru.NotifyUpdate(n)
  }
}
