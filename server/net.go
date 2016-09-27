package xtpctlserver

import (
  "io"
  "sync"

  manet "github.com/multiformats/go-multiaddr-net"

  pb "github.com/libp2p/go-xtp-ctl/pb"
)

type Transport struct {
  sync.RWMutex

  Id   int64
  Code string // multiaddr string code

  Listeners map[int64]*Listener
  Dialers   map[int64]*Dialer
  Conns     map[int64]*Conn

  NextId func() int64 // function to get the next id.
}

func NewTransport(nextId func() int64, code string) *Transport {
  return &Transport{
    Id:        nextId(), // allocate one for transport
    NextId:    nextId,
    Code:      code,
    Listeners: make(map[int64]*Listener),
    Dialers:   make(map[int64]*Dialer),
    Conns:     make(map[int64]*Conn),
  }
}

func (t *Transport) Close() error {
  t.Lock()
  defer t.Unlock()

  for id, l := range t.Listeners {
    l.Close()
    delete(t.Listeners, id)
  }

  for id, _ := range t.Dialers {
    delete(t.Dialers, id)
  }

  for id, c := range t.Conns {
    c.Close()
    delete(t.Conns, id)
  }

  return nil
}

func (t *Transport) CloseId(id int64) error {
  t.Lock()
  defer t.Unlock()

  l, found := t.Listeners[id]
  if found {
    delete(t.Listeners, id)
    return l.Close()
  }

  _, found = t.Dialers[id]
  if found {
    delete(t.Dialers, id)
    return nil
  }

  c, found := t.Conns[id]
  if found {
    delete(t.Conns, id)
    return c.Close()
  }

  for _, c := range t.Conns {
    c.CloseId(id)
  }
  return nil
}

type Listener struct {
  Id    int64
  Trans *Transport
  L     manet.Listener
}

func (t *Transport) NewListener(l manet.Listener) *Listener {
  id := t.NextId()
  l2 := &Listener{id, t, l}
  t.Lock()
  t.Listeners[id] = l2
  t.Unlock()
  return l2
}

func (l *Listener) Close() error {
  return l.L.Close()
}

type Dialer struct {
  Id    int64
  Trans *Transport
  D     *manet.Dialer
}

func (t *Transport) NewDialer(d *manet.Dialer) *Dialer {
  id := t.NextId()
  d2 := &Dialer{id, t, d}
  t.Lock()
  t.Dialers[id] = d2
  t.Unlock()
  return d2
}

type Conn struct {
  sync.RWMutex

  Id      int64
  C       manet.Conn

  Trans   *Transport
  Streams map[int64]*Stream
}

func (t *Transport) NewConn(mc manet.Conn) *Conn {
  id := t.NextId()
  c := &Conn{
    Id:      id,
    C:       mc,
    Trans:   t,
    Streams: make(map[int64]*Stream),
  }
  t.Lock()
  t.Conns[id] = c
  t.Unlock()
  return c
}

func (c *Conn) Close() error {
  c.Lock()
  for id, s := range c.Streams {
    s.Close()
    delete(c.Streams, id)
  }
  err := c.C.Close()
  c.Unlock()
  return err
}

func (c *Conn) CloseId(id int64) error {
  c.Lock()
  defer c.Unlock()
  s, found := c.Streams[id]
  if found {
    delete(c.Streams, id)
    return s.Close()
  }
  return nil
}

func (c *Conn) NewStream(st IoStream) *Stream {
  id := c.Trans.NextId()
  s := &Stream{id, st, c}
  c.Lock()
  c.Streams[id] = s
  c.Unlock()
  return s
}

type Stream struct {
  Id      int64
  S       IoStream
  Conn    *Conn
}

func (s *Stream) Close() error {
  return s.S.Close()
}

type IoStream interface {
  io.Reader
  io.Writer
  io.Closer
}

// protobuf helpers

func (t *Transport) PB() *pb.Transport {
  return &pb.Transport{
    Id:        &t.Id,
    Transport: &t.Code,
  }
}

func (l *Listener) PB() *pb.Listener {
  var b []byte
  if a := l.L.Multiaddr(); a != nil {
    b = a.Bytes()
  }

  return &pb.Listener{
    Id:          &l.Id,
    TransportId: &l.Trans.Id,
    Multiaddr:   b,
  }
}

func (d *Dialer) PB() *pb.Dialer {
  var b []byte
  if a := d.D.LocalAddr; a != nil {
    b = a.Bytes()
  }

  return &pb.Dialer{
    Id:          &d.Id,
    TransportId: &d.Trans.Id,
    Multiaddr:   b,
  }
}

func (c *Conn) PB() *pb.Conn {
  var lab, rab []byte
  if a := c.C.LocalMultiaddr(); a != nil {
    lab = a.Bytes()
  }
  if a := c.C.RemoteMultiaddr(); a != nil {
    rab = a.Bytes()
  }

  return &pb.Conn{
    Id:              &c.Id,
    TransportId:     &c.Trans.Id,
    LocalMultiaddr:  lab,
    RemoteMultiaddr: rab,
  }
}

func (s *Stream) PB() *pb.Stream {
  var lab, rab []byte
  if a := s.Conn.C.LocalMultiaddr(); a != nil {
    lab = a.Bytes()
  }
  if a := s.Conn.C.RemoteMultiaddr(); a != nil {
    rab = a.Bytes()
  }

  return &pb.Stream{
    Id:              &s.Id,
    ConnId:          &s.Conn.Id,
    TransportId:     &s.Conn.Trans.Id,
    LocalMultiaddr:  lab,
    RemoteMultiaddr: rab,
  }
}
