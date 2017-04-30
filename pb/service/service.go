// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package pbserv;

type Manager struct {
  worker.AbstractContext

}

func (m *Manager) Receive(m message.Envelope) {
  switch
  c := m.Division.WorkerWhence(ConnectionWorkerFactory())
  c.Send(m)
}

func (m *Manager) TcpListen() {
  go func() {
    defer m.listener.Close()
//    pbs.connection.(*net.TCPListener).SetDeadline(time.Now().Add(time.Second))
    for {
      select {
      case <- m.quitChannel:
          return
      default:
        conn, err := m.listener.Accept()
        if err != nil {
          log.Fatal("Error accepting: ", err.Error())
        }

        m.Send(&conn)

        nw, e := pbs.NewWorker(&conn)
        if e == nil {
          *pbs.workers = append(*pbs.workers, nw)
          nw.Start(pbs.task)
        }
      }
    }
  }()
}
