// Copyright 2017 Matthew Rich <matthewrich.conf@gmail.com>. All rights reserved.

package cluster;

import (
  "github.com/mrichconf/chatterbox/pb"
)

// store tag associations

type State struct {
  data map[string]proto.Message
}

func (s *State) Persist() {
  
}
