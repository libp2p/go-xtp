package xtpserver

import (
  "sync"
  "errors"

  xnet "github.com/libp2p/go-xtp-ctl/net"
)

type ServerClient struct {
  sync.RWMutex

  Server *Server
  Conn   xnet.Conn

  transports map[int64]*transport

  idCounter // embedded
}

// Close shuts down the ServerClient, closing everything.
func (sc *ServerClient) Close() error {
  panic("todo")
}

func (sc *ServerClient) transport(id int64) *transport {
  sc.Lock()
  t := sc.transports[id]
  sc.Unlock()
  return t
}

func (sc *ServerClient) addTransport(t *transport) {
  sc.Lock()
  sc.transports[t.id] = t
  sc.Unlock()
}

func (sc *ServerClient) rmTransport(t *transport) {
  sc.Lock()
  delete(sc.transports, t.id)
  sc.Unlock()
}

func (sc *ServerClient) Find(id int64) interface{} {
  sc.RLock()
  defer sc.RUnlock()
  t, found := sc.transports[id]
  if found {
    return t
  }

  for _, t := range sc.transports {
    v := t.Find(id)
    if v != nil {
      return v
    }
  }

  return nil
}

func (sc *ServerClient) CloseId(id int64) error {
  v := sc.Find(id)
  if v == nil {
    return nil // already closed?
  }

  switch v := v.(type) {
  case *transport:
    sc.rmTransport(v)
    return v.Close()
  case *listener:
    v.xport.rmListener(v)
    return v.Close()
  case *dialer:
    v.xport.rmDialer(v)
    return nil
  case *conn:
    v.xport.rmConn(v)
    return v.Close()
  case *stream:
    v.conn.rmStream(v)
    return v.Close()
  default:
    return errors.New("unknown type")
  }
}

type idCounter struct {
  lk   sync.Mutex
  next int64
}

func (c *idCounter) NextId() int64 {
  c.lk.Lock()
  c.next++
  n := c.next
  c.lk.Unlock()
  return n
}
