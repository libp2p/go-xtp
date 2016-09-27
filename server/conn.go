package xtpserver

import (
  "sync"

  xnet "github.com/libp2p/go-xtp-ctl/net"
)

type conn struct {
  sync.RWMutex

  id      int64
  rawC    xnet.Conn
  streams map[int64]*stream
  xport   *transport
}

func newConn(id int64, t *transport, c xnet.Conn) *conn {
  return &conn{
    id:      id,
    rawC:    c,
    xport:   t,
    streams: make(map[int64]*stream),
  }
}

func (c *conn) addStream(s *stream) {
  c.Lock()
  c.streams[s.id] = s
  c.Unlock()
}

func (c *conn) rmStream(s *stream) {
  c.Lock()
  delete(c.streams, s.id)
  c.Unlock()
}

func (c *conn) Close() error {
  c.Lock()
  for id, s := range c.streams {
    delete(c.streams, id)
    s.Close()
  }
  c.Unlock()
  return c.rawC.Close()
}

func (c *conn) Find(id int64) *stream {
  c.RLock()
  defer c.RUnlock()

  s, _ := c.streams[id]
  return s
}

func (c *conn) Dial() (*stream, error) {
  s, err := c.rawC.Dial()
  if err != nil {
    return nil, err
  }
  id := c.xport.sc.NextId()

  s2 := newStream(id, c, s)
  c.addStream(s2)
  return s2, nil
}

func (c *conn) Accept() (*stream, error) {
  s, err := c.rawC.Accept()
  if err != nil {
    return nil, err
  }
  id := c.xport.sc.NextId()

  s2 := newStream(id, c, s)
  c.addStream(s2)
  return s2, nil
}
