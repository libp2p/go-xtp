package xtpctlclient

import (
  ma "github.com/multiformats/go-multiaddr"
  manet "github.com/multiformats/go-multiaddr-net"
  xnet "github.com/libp2p/go-xtp-ctl/net"
  xrpc "github.com/libp2p/go-xtp-ctl/rpc"
)

type listener struct {
  id     int64
  tid    int64
  ctls   IoStream     // the xtp-ctl stream for this listener.
  client *Client      // the xtp-ctl client
  laddr  ma.Multiaddr // the address of this listener.
}

// Multiaddr returns the listener's (local) Multiaddr.
func (l *listener) Multiaddr() ma.Multiaddr {
  return l.laddr
}

// Accept waits for and returns the next connection to the listener.
// Returns a Multiaddr friendly Conn
func (l *listener) Accept() (xnet.Conn, error) {
  // open a new data stream
  s, err := l.client.Conn.Dial()
  if err != nil {
    return nil, err
  }

  // Send an accept request, wait for an accept response
  res, err := xrpc.AcceptReq(s, l.id, raddr)
  if err != nil {
    return err
  }

  c := newConn(l.client.Conn, s, res.Conn)
  return c, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *listener) Close() error {
  _, err := xrpc.CloseReq(l.ctls, l.id)
  l.ctls.Close()
  return err
}

func newListener(c *Client, ctls xnet.Stream, l *pb.Listener) (*listener, error) {
  if !l.Valid() {
    return nil, xrpc.ErrInvalidResponse
  }
  a, err := ma.NewMultiaddrBytes(l.Multiaddr)
  if err != nil {
    return nil, err
  }

  return &listener{
    id:     *l.Id,
    tid:    *l.TransportId,
    ctls:   ctls,
    client: c,
    laddr:  a,
  }, nil
}
