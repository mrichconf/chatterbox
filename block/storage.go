// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package block;

type BuildIndex {}

type Storage struct {
  worker.AbstractContext
  index map[[32]byte]int64
  io fs.IFile
}


func (t *Storage) Receive(m message.Envelope) {
  s := m.Sender()
  switch b := m.Body().(type) {
  case *BuildIndex:
    fi,e := t.io.Stat()
    if e != nil {
      log.Fatal(e)
    }
    fSz := fi.Size()

    if fSz == 0 {
      binary.Write(t.io, binary.LittleEndian, int32(t.blockSize))
    } else {
      var size int32
      se := binary.Read(t.io, binary.LittleEndian, &size)
      if se != nil {
        log.Fatal(se)
      }
      if size != int32(t.blockSize) {
        log.Fatal("Requested block size (", t.blockSize, ") doesn't match datastore (", size, ")")
      }
    }

    partialBlock := (fSz - 4) % int64(t.blockSize)
    if partialBlock != 0 {
      // TODO: it might be better to store the block Ids and verify the data
      log.Println("Data file contains a partial block - truncating ", partialBlock, " bytes from end.")
      t.io.Truncate(fSz - partialBlock)
    }

    data := make([]byte, t.blockSize)
    for {
      offset, _ := t.io.Seek(0,1)
      n, e := t.io.Read(data)
      if e != nil {
        // reading error
      }
      if n == 0 {
        break
      }
      if n < t.blockSize {
        // corrupted?
        log.Fatal("Unable to read complete block. Bytes read:", n)
      }
      t.index[sha512.Sum512_256(data)] = offset
    }

  case *Block:
    // ensure this block is persisted 
    response := &Block{ Size: proto.Int32(int32(t.blockSize)) }

    id,e := b.GenId()
    if e != nil {
      response.State = State_absent.Enum()
    }

    i, exists := t.index[id]
    
    if b.Data == nil || len(b.Data) != t.blockSize {
      if exists {
        if b.Flags == Flags_excludeData.Enum() {
        } else {
          data = make([]byte, t.blockSize)
          n, e := t.io.ReadAt(data, i)
          if e != nil {
            // read error or io.EOF
            log.Fatal("Failed reading at ", int64(4 + (int64(b.blockSize) * (i - 1))), " bytes: ", n, " err", e)
          }
          if n != b.blockSize {
            log.Fatal("Failed reading block from disk.  Bytes read: ", n)
          }

          response.Data = data
        }
        response.State = State_present.Enum()
      } else {
        response.State = State_absent.Enum()
      }
    } else {
      if exists {
        // but the block alreadys exists: do nothing
      } else {
        b.index[id], _ = b.file.Seek(0, 1)

        // write block to fs
        b.file.Write(data)
      }
      response.Block.Flags = Flags_excludeData.Enum()
      response.Block.Data = nil
      response.Block.State = State_present.Enum()
    }
    response.Id = id[:]
    s.Send(message.New(response, t))
  }
}
