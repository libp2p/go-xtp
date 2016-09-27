package xtp_ctl

// this file includes some helper functions added to the protobuf structs

import proto "github.com/gogo/protobuf/proto"
import ma "github.com/multiformats/go-multiaddr"

const MinId = 1

type ListReqTypes struct {
  Transports bool
  Listeners  bool
  Dialers    bool
  Conns      bool
  Streams    bool
}

func (m *ListReq) TypesRequested() (t ListReqTypes) {
  for typ := range m.Types {
    switch TType(typ) {
    case TType_TTypeNil: // ignore
    case TType_TTypeTransport:
      t.Transports = true
    case TType_TTypeListener:
      t.Listeners = true
    case TType_TTypeDialer:
      t.Dialers = true
    case TType_TTypeConn:
      t.Conns = true
    case TType_TTypeStream:
      t.Streams = true
    }
  }
  return t
}

func Mk_ListRes_Item(id int64, typ TType, val proto.Message) (*ListRes_Item, error) {
  buf, err := proto.Marshal(val)
  if err != nil {
    return nil, err
  }
  return &ListRes_Item{
    Id:    &id,
    Type:  &typ,
    Value: buf,
  }, nil
}

func ListRes_Item_Transport(t *Transport) (*ListRes_Item, error) {
  return Mk_ListRes_Item(*t.Id, TType_TTypeTransport, t)
}

func ListRes_Item_Listener(l *Listener) (*ListRes_Item, error) {
  return Mk_ListRes_Item(*l.Id, TType_TTypeListener, l)
}

func ListRes_Item_Dialer(d *Dialer) (*ListRes_Item, error) {
  return Mk_ListRes_Item(*d.Id, TType_TTypeDialer, d)
}

func ListRes_Item_Conn(c *Conn) (*ListRes_Item, error) {
  return Mk_ListRes_Item(*c.Id, TType_TTypeConn, c)
}

func ListRes_Item_Stream(s *Stream) (*ListRes_Item, error) {
  return Mk_ListRes_Item(*s.Id, TType_TTypeStream, s)
}

// validators

func (m *RPC) Valid() bool {
  if m == nil || m.Rpc == nil {
    return false
  }
  return true
}

func (m *ListRes_Item) Valid() bool {
  if m == nil || m.Id == nil || m.Type == nil || m.Value == nil {
    return false
  }
  if *m.Id < 1 {
    return false
  }
  return true
}

func (m *ListenRes) Valid() bool {
  if m == nil || m.Listener == nil {
    return false
  }
  return m.Listener.Valid()
}

func (m *DialerRes) Valid() bool {
  if m == nil || m.Dialer == nil {
    return false
  }
  return m.Dialer.Valid()
}

func (m *Transport) Valid() bool {
  if m == nil || m.Id == nil || m.Transport == nil {
    return false
  }
  if *m.Id <= MinId || *m.Transport == "" {
    return false
  }
  return true
}

func (m *Listener) Valid() bool {
  if m == nil || m.Id == nil || m.TransportId == nil || m.Multiaddr == nil {
    return false
  }
  if *m.Id < 1 || *m.TransportId < 1 {
    return false
  }
  if _, err := ma.NewMultiaddrBytes(m.Multiaddr); err != nil {
    return false
  }
  return true
}

func (m *Dialer) Valid() bool {
  if m == nil || m.Id == nil || m.TransportId == nil || m.Multiaddr == nil {
    return false
  }
  if *m.Id < 1 || *m.TransportId < 1 {
    return false
  }
  if _, err := ma.NewMultiaddrBytes(m.Multiaddr); err != nil {
    return false
  }
  return true
}

func (m *Conn) Valid() bool {
  if m == nil || m.Id == nil || m.TransportId == nil || m.LocalMultiaddr == nil || m.RemoteMultiaddr == nil {
    return false
  }
  if *m.Id < 1 || *m.TransportId < 1 {
    return false
  }
  if _, err := ma.NewMultiaddrBytes(m.LocalMultiaddr); err != nil {
    return false
  }
  if _, err := ma.NewMultiaddrBytes(m.RemoteMultiaddr); err != nil {
    return false
  }
  return true
}

func (m *Stream) Valid() bool {
  if m == nil || m.Id == nil || m.TransportId == nil || m.ConnId == nil ||
    m.LocalMultiaddr == nil || m.RemoteMultiaddr == nil {
    return false
  }
  if *m.Id < 1 || *m.TransportId < 1 || *m.ConnId < 1 {
    return false
  }
  if _, err := ma.NewMultiaddrBytes(m.LocalMultiaddr); err != nil {
    return false
  }
  if _, err := ma.NewMultiaddrBytes(m.RemoteMultiaddr); err != nil {
    return false
  }
  return true
}
