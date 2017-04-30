// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package pb;


// codec messages
// TODO: generic codec messages
type Encode struct {
  Message interface{}
}
type Decode struct {
  Message interface{}
}

// enc/dec protobuf messages 
type CodecWorker struct {
  AbstractWorker
  client worker.Ref // nio.Client
  messageFactory worker.Factory
}

func (c *Client) Receive(m message.Envelope) {
  s := m.Sender()
  switch b := m.Body().(type) {
  case *message.Register:
    if b.Handler == nil {
      //t.target = b.Sender()
      //error
    } else {
      t.messageFactory = b.Handler.(worker.Factory)
    }
  case *ReadResponse: // response from the connection
    pb := t.messageFactory()
    err = proto.Unmarshal(b.Content, pb)
    if err != nil {
      return err
    }
    b.Origin.Send(message.New(pb, c))
  case *WriteResponse:
    
  case *Encode:
    data, err := proto.Marshal(b.Message.(proto.Message))
    if err != nil {
      return 0, err
    }
    z := make([]byte, 4)
    binary.LittleEndian.PutUint32(z, uint32(len(data)))
    d = append(z, data...)
    c.client.Send(message.New(&WriteRequest{ Content: d, Origin: s }, c))
  case *Decode:
    var size uint32
    // read next message size
    i := NewInterrogator()
    c.client.Send(message.New(&ReadRequest{ Size: 4 }, i))
    m := <- i
    size = binary.LittleEndian.Uint32(m.Body().([]byte))

    c.client.Send(message.New(&ReadRequest{ Size: size, Origin: s }, c)) // forward to our client worker
  }
}
