package xtpctlclient

import (
  ma "github.com/multiformats/go-multiaddr"
  manet "github.com/multiformats/go-multiaddr-net"
  xnet "github.com/libp2p/go-xtp-ctl/net"
  xrpc "github.com/libp2p/go-xtp-ctl/rpc"
)

type stream struct {
  id     int64
  tid    int64
  ctls   IoStream // the xtp-ctl stream for this conn
  client *Client  // the xtp-ctl client
  conn   *Conn    // the conn this stream belongs to
}


// LocalMultiaddr returns the local Multiaddr associated
// with this connection
func (c *conn) LocalMultiaddr() ma.Multiaddr {
  return c.laddr
}

// RemoteMultiaddr returns the remote Multiaddr associated
// with this connection
func (c *conn) RemoteMultiaddr() ma.Multiaddr {
  return c.raddr
}

// Close closes the dialer.
func (s *stream) Close() error {
  _, err := xrpc.CloseReq(s.ctls, s.id)
  s.ctls.Close()
  return err
}

func newStream(c *Client, ctls xnet.Stream, s *pb.Stream, cn *conn) (*dialer, error) {
  if !s.Valid() {
    return nil, xrpc.ErrInvalidResponse
  }
  return &stream{
    id:     *s.Id,
    tid:    *s.TransportId,
    ctls:   ctls,
    client: c,
    conn:   cn,
  }, nil
}

type IoStream interface {
  io.Reader
  io.Writer
  io.Closer
}
