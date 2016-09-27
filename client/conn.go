package xtpctlclient

import (
  pb "github.com/libp2p/go-xtp-ctl/pb"
  ma "github.com/multiformats/go-multiaddr"
  xnet "github.com/libp2p/go-xtp-ctl/net"
  xrpc "github.com/libp2p/go-xtp-ctl/rpc"
)

type conn struct {
  id     int64
  tid    int64
  ctls   IoStream     // the xtp-ctl stream for this conn
  client *Client      // the xtp-ctl client
  laddr  ma.Multiaddr // the local address of this conn
  raddr  ma.Multiaddr // the remote address of this conn
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

// Dial attempts to open a new stream across Conn to the other side.
func (c *conn) Dial() (xnet.Stream, error) {
  // open a new data stream
  s, err := c.client.Conn.Dial()
  if err != nil {
    return nil, err
  }

  // Send an accept request, wait for an accept response
  res, err := xrpc.DialReq(s, c.id, nil)
  if err != nil {
    return nil, err
  }

  return newStream(c.client, s, res.Stream, c)
}

// Accept accepts an incoming conn.Dial from the other side.
func (c *conn) Accept() (xnet.Stream, error) {
  // open a new data stream
  s, err := c.client.Conn.Dial()
  if err != nil {
    return nil, err
  }

  // Send an accept request, wait for an accept response
  res, err := xrpc.AcceptReq(s, c.id)
  if err != nil {
    return nil, err
  }

  return newStream(c.client, s, res.Stream, c)
}

// Close closes the dialer.
func (c *conn) Close() error {
  err := xrpc.CloseReq(c.ctls, c.id)
  c.ctls.Close()
  return err
}

func newConn(c *Client, ctls xnet.Stream, cn *pb.Conn) (*conn, error) {
  if !cn.Valid() {
    return nil, xrpc.ErrInvalidResponse
  }
  lm, err := ma.NewMultiaddrBytes(cn.LocalMultiaddr)
  if err != nil {
    return nil, err
  }
  rm, err := ma.NewMultiaddrBytes(cn.RemoteMultiaddr)
  if err != nil {
    return nil, err
  }

  return &conn{
    id:     *cn.Id,
    tid:    *cn.TransportId,
    ctls:   ctls,
    client: c,
    laddr:  lm,
    raddr:  rm,
  }, nil
}
