// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package cluster

import (
  "log"
  "sync"
  "github.com/mrichconf/chatterbox/cluster/proto"
  "github.com/mrichconf/feudal/worker"
  "github.com/golang/protobuf/proto"
//  "github.com/golang/groupcache/groupcachepb"
)

type QueryRequestInjector struct {
  messenger ProtobufMessenger
  task worker.Task
  messageTasks [255]worker.Task
}
type QueryResponseInjector struct {
  sync.RWMutex
  responses map[uint32]chan proto.Message
}

func NewQueryRequestInjector(t worker.Task) *QueryRequestInjector {
  return &QueryRequestInjector{ task: t }
}
func NewQueryResponseInjector() *QueryResponseInjector {
  return &QueryResponseInjector{ responses: make(map[uint32]chan proto.Message, 50) }
}

func (qi *QueryRequestInjector) SetMessageTask(mType uint32, task worker.Task) {
  qi.messageTasks[mType] = task
}

func (qi *QueryRequestInjector) NewMessage(r ProtobufReceiver) error {
  m := &cluster_proto.QueryRequest{}
  if e := r.Receive(m); e != nil {
    return e
  }
  if m.MessageType != cluster_proto.Type_query {
    return nil
  }
  log.Println("QueryRequestInjector::NewMessage(): ", m)

//  g,_ := qi.messenger.DecodeMessage(m.Query.Body.Type, m.Query.Body.Value)

  var err error
  var resp proto.Message
  log.Println("QueryRequestInjector::NewMessage(): - looking up messageTask for type: ", m.Query.Body.Type)
  if qi.messageTasks[m.Query.Body.Type] != nil {
    resp, err = qi.messageTasks[m.Query.Body.Type].Handler(nil, m)
  } else {
    // does this instance have the requested block?
    resp, err = qi.task.Handler(nil, m)
  }
  if err != nil {
    return err
  }
  if resp != nil {
    rQ := resp.(*cluster_proto.QueryResponse)
    if rQ.MessageType == cluster_proto.Type_response {
      qi.messenger.Queue(resp)
    }
  }
  return nil
}

func (qi *QueryRequestInjector) Receive(msg proto.Message) error {
/*
  r := msg.(*cluster_proto.QueryRequest)
  id := r.Query.Id
  qi.RLock()
  defer qi.RUnlock()
  if _,exists := qi.responses[id]; !exists {
    qi.Lock()
    qi.responses[id] = make(chan proto.Message, 50)
    qi.Unlock()
  }
  select {
  case rM := <- qi.responses[id]:
    msg = rM
  }
*/
  return nil
}

// handle a queryresponse
func (qi *QueryResponseInjector) NewMessage(r ProtobufReceiver) error {
  m := &cluster_proto.QueryResponse{}
  if e := r.Receive(m); e != nil {
    return e
  }
  if m.MessageType != cluster_proto.Type_response {
    return nil
  }
  log.Println("QueryResponseInjector::NewMessage(): ", m)
  qi.Lock()
  defer qi.Unlock()
  id := m.Response.Id
  if _,exists := qi.responses[id]; !exists {
    log.Println("QueryResponseInjector::NewMessage(): creating response channel: ", id)
    qi.responses[id] = make(chan proto.Message, 50)
  }
  log.Println("QueryResponseInjector::NewMessage(): writing message to channel")
  qi.responses[id] <- m
  return nil
}

func (qi *QueryResponseInjector) Receive(msg proto.Message) error {
  r := msg.(*cluster_proto.QueryResponse)
  id := r.Response.Id
  qi.Lock()
  if _,exists := qi.responses[id]; !exists {
    log.Println("QueryResponseInjector::Receive() creating response channel: ", id)
    qi.responses[id] = make(chan proto.Message, 50)
  }
  ch := qi.responses[id]
  qi.Unlock()
  log.Println("waiting for response: ", id)
  select {
  case rM := <- ch:
    log.Println("QueryResponseInjector::Receive() received response: ", rM)
    ir := rM.(*cluster_proto.QueryResponse)
    r.Response = ir.Response
    log.Println("QueryResponseInjector::Receive() received response: ", r)
    //msg = rM
  }
  return nil
}
