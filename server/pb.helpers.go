package xtpserver

import (
  pb "github.com/libp2p/go-xtp-ctl/pb"
)

// protobuf helpers

func (t *transport) PB() *pb.Transport {
  code := t.rawT.Code()
  return &pb.Transport{
    Id:        &t.id,
    Transport: &code,
  }
}

func (l *listener) PB() *pb.Listener {
  var b []byte
  if a := l.rawL.Multiaddr(); a != nil {
    b = a.Bytes()
  }

  return &pb.Listener{
    Id:          &l.id,
    TransportId: &l.xport.id,
    Multiaddr:   b,
  }
}

func (d *dialer) PB() *pb.Dialer {
  var b []byte
  if a := d.rawD.Multiaddr(); a != nil {
    b = a.Bytes()
  }

  return &pb.Dialer{
    Id:          &d.id,
    TransportId: &d.xport.id,
    Multiaddr:   b,
  }
}

func (c *conn) PB() *pb.Conn {
  var lab, rab []byte
  if a := c.rawC.LocalMultiaddr(); a != nil {
    lab = a.Bytes()
  }
  if a := c.rawC.RemoteMultiaddr(); a != nil {
    rab = a.Bytes()
  }

  return &pb.Conn{
    Id:              &c.id,
    TransportId:     &c.xport.id,
    LocalMultiaddr:  lab,
    RemoteMultiaddr: rab,
  }
}

func (s *stream) PB() *pb.Stream {
  var lab, rab []byte
  if a := s.conn.rawC.LocalMultiaddr(); a != nil {
    lab = a.Bytes()
  }
  if a := s.conn.rawC.RemoteMultiaddr(); a != nil {
    rab = a.Bytes()
  }

  return &pb.Stream{
    Id:              &s.id,
    ConnId:          &s.conn.id,
    TransportId:     &s.conn.xport.id,
    LocalMultiaddr:  lab,
    RemoteMultiaddr: rab,
  }
}
