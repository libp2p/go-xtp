package xtpctlclient

import (
  xnet "github.com/libp2p/go-xtp-ctl/net"
)

type Client struct {
  Conn xnet.Conn // connection to the server
}
