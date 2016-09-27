package xtpctlclient

import (
  ma "github.com/multiformats/go-multiaddr"
  manet "github.com/multiformats/go-multiaddr-net"
  xnet "github.com/libp2p/go-xtp-ctl/net"
  xrpc "github.com/libp2p/go-xtp-ctl/rpc"
)

type dialer struct {
  id     int64
  tid    int64
  ctls   IoStream     // the xtp-ctl stream for this dialer
  client *Client      // the xtp-ctl client
  laddr  ma.Multiaddr // the address of this dialer
}

// Multiaddr returns the dialer's (local) Multiaddr.
func (d *dialer) Multiaddr() ma.Multiaddr {
  return d.laddr
}

// Dial dials the given multiaddr and sets up a connection.
func (d *dialer) Dial(raddr ma.Multiaddr) (xnet.Conn, error) {
  // open a new data stream
  s, err := d.client.Conn.Dial()
  if err != nil {
    return nil, err
  }

  // Send an accept request, wait for an accept response
  res, err := xrpc.DialReq(s, d.id, raddr)
  if err != nil {
    return err
  }

  c := newConn(l.client.Conn, s, res.Conn)
  return c, nil
}

// Close closes the dialer.
func (d *dialer) Close() error {
  _, err := xrpc.CloseReq(d.ctls, d.id)
  d.ctls.Close()
  return err
}

func newDialer(c *Client, ctls xnet.Stream, d *pb.Dialer) (*dialer, error) {
  if !d.Valid() {
    return nil, xrpc.ErrInvalidResponse
  }
  a, err := ma.NewMultiaddrBytes(d.Multiaddr)
  if err != nil {
    return nil, err
  }

  return &dialer{
    id:     *d.Id,
    tid:    *d.TransportId,
    ctls:   ctls,
    client: c,
    laddr:  a,
  }, nil
}
