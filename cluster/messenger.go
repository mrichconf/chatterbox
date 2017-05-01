// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package cluster

import (
  "log"
  "errors"
  "github.com/hashicorp/memberlist"
  "github.com/golang/protobuf/proto"
  "github.com/mrichconf/feudal/message"
)

type ReceivedMessage struct {
  data []byte
}

type bc struct {
  message []byte
  finishedChannel chan bool
}

type Messenger struct {
  broadcastQueue *memberlist.TransmitLimitedQueue
  receiveChannel chan []byte
  messageTypes [255]message.Factory
  messageInjectors []MessageInjector
}

func NewMessenger(numNodes func() int, retransmitmult int) *Messenger {
  return &Messenger{
    broadcastQueue: &memberlist.TransmitLimitedQueue{ NumNodes: numNodes, RetransmitMult: retransmitmult }, 
    receiveChannel: make(chan []byte, 50),
  }
}

func (rm *ReceivedMessage) Receive(msg proto.Message) error {
  return proto.Unmarshal(rm.data, msg)
}

// FIXME - handle shutdown
func (cm *Messenger) Receiver() {
  log.Println("Starting cluster message receiver")
  go func() {
    for {
      select {
      case m := <- cm.receiveChannel:
        // call registered MessageType
        for i := range cm.messageInjectors {
          log.Println("Trying message type: ", i)
          go cm.messageInjectors[i].NewMessage(&ReceivedMessage{ data: m })
        }
      }
    }
  }()
}

func (cm *Messenger) RegisterMessageInjector(i MessageInjector) {
  cm.messageInjectors = append(cm.messageInjectors, i)
}

func (cm *Messenger) RegisterMessageType(mType uint32, m message.Factory) error {
  log.Println("RegisterMessageType() ", mType, m)
  if mType > 255 {
    return errors.New("Type identifier out of bounds")
  }
  cm.messageTypes[mType] = m
  return nil
}

func (cm *Messenger) DecodeMessage(mType uint32, data []byte) (proto.Message, error) {
  log.Println("DecodeMessage() ", mType, cm.messageTypes[mType]) 
  pb := cm.messageTypes[mType].NewMessage()
  e := proto.Unmarshal(data, pb)
  return pb, e
}

func (cm *Messenger) Len() int {
  return cm.broadcastQueue.NumQueued()
}

func (cm *Messenger) NodeMeta(limit int) []byte {
  return nil
}

func (cm *Messenger) NotifyMsg(buf []byte) {
  log.Println("NotifyMsg() ", buf)
  if len(buf) > 0 {
    if len(cm.receiveChannel) == cap(cm.receiveChannel) {
      panic("receive channel full")
    } else {
      cm.receiveChannel <- buf
    }
  }
}

func (cm *Messenger) GetBroadcasts(overhead, limit int) [][]byte {
  r := cm.broadcastQueue.GetBroadcasts(overhead, limit)
  if len(r) > 0 {
    log.Println("Messenger::GetBroadcasts() ", r)
  }
  return r
}

func (cm *Messenger) LocalState(join bool) []byte {
  return nil
}

func (cm *Messenger) MergeRemoteState(buf []byte, join bool) {

}

func (cm *Messenger) Queue(msg proto.Message) error {
  log.Println("Queue() ")
  d, e := proto.Marshal(msg)
  if e != nil {
    return e
  }
  cm.broadcastQueue.QueueBroadcast(&bc{ message: d })
  return nil
}

func (cm *Messenger) Receive(msg proto.Message) error {
  log.Println("Receive() ", len(cm.receiveChannel), "in receive channel")
  buf := <- cm.receiveChannel
  log.Println("Receive() Unmarshaling", buf)
  return proto.Unmarshal(buf, msg)
}

func (b *bc) Invalidates(br memberlist.Broadcast) bool {
  return false
}

func (b *bc) Message() []byte {
  return b.message
}

func (b *bc) Finished() {
  if b.finishedChannel != nil {
    close(b.finishedChannel)
  }
}
