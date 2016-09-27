package xtpserver

import (
  xnet "github.com/libp2p/go-xtp-ctl/net"
  ma "github.com/multiformats/go-multiaddr"
)

type dialer struct {
  id    int64
  rawD  xnet.Dialer
  xport *transport
}

func newDialer(id int64, t *transport, d xnet.Dialer) *dialer {
  return &dialer{id, d, t}
}


func (d *dialer) Dial(raddr ma.Multiaddr) (*conn, error) {
  c, err := d.rawD.Dial(raddr)
  if err != nil {
    return nil, err
  }
  id := d.xport.sc.NextId()

  c2 := newConn(id, d.xport, c)
  d.xport.addConn(c2)
  return c2, nil
}
