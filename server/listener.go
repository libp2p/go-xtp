package xtpserver

import (
  xnet "github.com/libp2p/go-xtp-ctl/net"
)

type listener struct {
  id    int64
  rawL  xnet.Listener
  xport *transport
}

func newListener(id int64, t *transport, l xnet.Listener) *listener {
  return &listener{id, l, t}
}

func (l *listener) Accept() (*conn, error) {
  c, err := l.rawL.Accept()
  if err != nil {
    return nil, err
  }
  id := l.xport.sc.NextId()

  c2 := newConn(id, l.xport, c)
  l.xport.addConn(c2)
  return c2, nil
}

func (l *listener) Close() error {
  return l.rawL.Close()
}
