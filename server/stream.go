package xtpserver

import (
  xnet "github.com/libp2p/go-xtp-ctl/net"
)

type stream struct {
  id    int64
  rawS  xnet.Stream
  conn  *conn
}

func newStream(id int64, c *conn, s xnet.Stream) *stream {
  return &stream{id, s, c}
}

func (s *stream) Close() error {
  return s.rawS.Close()
}
