package xtpserver

import (
  "sync"

  ma "github.com/multiformats/go-multiaddr"
  pb "github.com/libp2p/go-xtp-ctl/pb"
  xnet "github.com/libp2p/go-xtp-ctl/net"
)

type transport struct {
  sync.RWMutex

  id   int64
  rawT xnet.Transport
  sc   *ServerClient

  listeners map[int64]*listener
  dialers   map[int64]*dialer
  conns     map[int64]*conn
}

func newTransport(id int64, sc *ServerClient, t xnet.Transport) *transport {
  return &transport{
    id:   id,
    rawT: t,
    sc:   sc,

    listeners: make(map[int64]*listener),
    dialers:   make(map[int64]*dialer),
    conns:     make(map[int64]*conn),
  }
}

func (t *transport) listener(id int64) *listener {
  t.Lock()
  l := t.listeners[id]
  t.Unlock()
  return l
}

func (t *transport) addListener(l *listener) {
  t.Lock()
  t.listeners[l.id] = l
  t.Unlock()
}

func (t *transport) rmListener(l *listener) {
  t.Lock()
  delete(t.listeners, l.id)
  t.Unlock()
}

func (t *transport) dialer(id int64) *dialer {
  t.Lock()
  d := t.dialers[id]
  t.Unlock()
  return d
}

func (t *transport) addDialer(d *dialer) {
  t.Lock()
  t.dialers[d.id] = d
  t.Unlock()
}

func (t *transport) rmDialer(d *dialer) {
  t.Lock()
  delete(t.dialers, d.id)
  t.Unlock()
}

func (t *transport) conn(id int64) *conn {
  t.Lock()
  c := t.conns[id]
  t.Unlock()
  return c
}

func (t *transport) addConn(c *conn) {
  t.Lock()
  t.conns[c.id] = c
  t.Unlock()
}

func (t *transport) rmConn(c *conn) {
  t.Lock()
  delete(t.conns, c.id)
  t.Unlock()
}

func (t *transport) Close() error {
  t.Lock()
  defer t.Unlock()

  for id, l := range t.listeners {
    l.Close()
    delete(t.listeners, id)
  }

  for id, _ := range t.dialers {
    delete(t.dialers, id)
  }

  for id, c := range t.conns {
    delete(t.conns, id)
    c.Close()
  }

  return nil
}

func (t *transport) Find(id int64) interface{} {
  t.RLock()
  defer t.RUnlock()

  if t.id == id {
    return t
  }

  if l, ok := t.listeners[id]; ok {
    return l
  }

  if d, ok := t.dialers[id]; ok {
    return d
  }

  if c, ok := t.conns[id]; ok {
    return c
  }

  for _, c := range t.conns {
    s := c.Find(id)
    if s != nil {
      return s
    }
  }
  return nil
}

func (t *transport) List(types pb.ListReqTypes) []*pb.ListRes_Item {
  t.RLock()
  defer t.Unlock()

  var items []*pb.ListRes_Item

  addItem := func(i *pb.ListRes_Item, err error) {
    if err != nil {
      items = append(items, i)
    } else {
      // internal error. TODO: log it
    }
  }

  if types.Listeners {
    for _, l := range t.listeners {
      i, err := pb.ListRes_Item_Listener(l.PB())
      addItem(i, err)
    }
  }

  if types.Dialers {
    for _, d := range t.dialers {
      i, err := pb.ListRes_Item_Dialer(d.PB())
      addItem(i, err)
    }
  }

  if types.Conns {
    for _, c := range t.conns {
      i, err := pb.ListRes_Item_Conn(c.PB())
      addItem(i, err)
    }
  }

  if types.Streams {
    for _, c := range t.conns {
      c.Lock()
      for _, s := range c.streams {
        i, err := pb.ListRes_Item_Stream(s.PB())
        addItem(i, err)
      }
      c.Unlock()
    }
  }
  return items
}

func (t *transport) Listen(laddr ma.Multiaddr) (*listener, error) {
  l, err := t.rawT.Listen(laddr)
  if err != nil {
    return nil, err
  }
  id := t.sc.NextId()

  l2 := newListener(id, t, l)
  t.addListener(l2)
  return l2, nil
}

func (t *transport) Dialer(laddr ma.Multiaddr) (*dialer, error) {
  d, err := t.rawT.Dialer(laddr)
  if err != nil {
    return nil, err
  }
  id := t.sc.NextId()

  d2 := newDialer(id, t, d)
  t.addDialer(d2)
  return d2, nil
}

func (t *transport) Dial(raddr ma.Multiaddr) (*conn, error) {
  c, err := t.rawT.Dial(raddr)
  if err != nil {
    return nil, err
  }
  id := t.sc.NextId()

  c2 := newConn(id, t, c)
  t.addConn(c2)
  return c2, nil
}
