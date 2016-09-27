package xtpctlserver

import (
  "errors"

  manet "github.com/multiformats/go-multiaddr-net"
)

type Server struct {
  Listeners []manet.Listener
  Clients   []ServerClient

  NextId func() int64 // function to get the next id.
}
