// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package pbserv;

// Runner to dispatch protobuf messages to connection worker
type Protocol struct {
  nio.Connection
  quit chan bool
  stopped chan bool
}

// Reads the next protobuf message found in the file stream
func (p *Protocol) Start(d worker.Dispatcher, c worker.Context) {
  go func() {
    for {
      select {
      case <- p.quit:
        p.Connection.Close()
        if v := c.Division(); v != nil {
          v.Disolve()
        }
        close(s.stopped)
        return
      default:
        var size int32
        if e := binary.Read(p.Connection, binary.LittleEndian, &size); e != nil || size < 1 {
          c.Manager().Send(message.New(&message.Error{ e: e }, c)) //send an interrogator and wait for the response?
          // should we close the connection? or just ignore it
          panic(e)
        }

        buffer := make([]byte, size)
        n, err := io.ReadAtLeast(p.connection, buffer, int(size))
        if err == io.EOF {
          // we're done right?
          return err
        }
        if int32(n) != size {
          // error
          return err
        }
        if err != nil {
          // io error
          return err
        }
        d.Dispatch(message.New(&message.Serialized{ content: buffer }, c))
      }
    }
  }()
}
