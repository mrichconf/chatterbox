// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package blockproto;

func (b *Block) GenId() (id [32]byte, err error) {
  data := b.GetData()
  blockId := b.GetId()
  if blockId != nil {
    copy(id[:], blockId)
  } else if data != nil {
    id = sha512.Sum512_256(data)
  } else {
    err = errors.New("Unable to get block id")
  }
  return
}
