package xtpclient

import (
  pb "github.com/libp2p/go-xtp-ctl/pb"
  ma "github.com/multiformats/go-multiaddr"
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
  res, err := xrpc.AcceptReq(s, l.id)
  if err != nil {
    return nil, err
  }

  return newConn(l.client, s, res.Conn)
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *listener) Close() error {
  err := xrpc.CloseReq(l.ctls, l.id)
  l.ctls.Close()
  return err
}

func newListener(c *Client, ctls xnet.Stream, l *pb.Listener) (*listener, error) {
  if !l.Valid() {
    return nil, xrpc.ErrInvalidMessage
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
