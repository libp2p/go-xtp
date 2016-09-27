package xtpctlclient

import (
  ma "github.com/multiformats/go-multiaddr"
  manet "github.com/multiformats/go-multiaddr-net"
  xnet "github.com/libp2p/go-xtp-ctl/net"
  xrpc "github.com/libp2p/go-xtp-ctl/rpc"
  pb "github.com/libp2p/go-xtp-ctl/pb"
)

type transport struct {
  id     int64
  ctls   IoStream // the xtp-ctl stream for this listener.
  client *Client
  code   string
}

func (t *transport) Code() string {
  return t.code
}

func (t *transport) Listen(raddr ma.Multiaddr) (xnet.Listener, error) {
  // open a new control stream
  s, err := t.client.Conn.Dial()
  if err != nil {
    return nil, err
  }

  // Send a listen request, wait for a listen response
  res, err := xrpc.ListenReq(s, t.id, raddr)
  if err != nil {
    return err
  }

  return newListener(t.client, s, res.Listener)
}

func (t *transport) Dial(raddr ma.Multiaddr) (xnet.Conn, error) {
  // open a new data stream
  s, err := t.client.Conn.Dial()
  if err != nil {
    return nil, err
  }

  // Send a dial request, wait for the dial response
  res, err := xrpc.DialReq(s, t.id, raddr)
  if err != nil {
    return err
  }

  c := newConn(t.client, s, res.Conn)
  return c, nil
}

func (t *transport) Dialer(laddr ma.Multiaddr) (xnet.Dialer, error) {
  // open a new control stream
  s, err := t.client.Conn.Dial()
  if err != nil {
    return nil, err
  }

  // Send a dialer request, wait for the dialer response
  res, err := xrpc.DialerReq(s, t.id, laddr)
  if err != nil {
    return err
  }

  c := newDialer(t.client, s, res.Dialer)
  return c, nil
}

func (t *transport) Close() error {
  _, err := xrpc.CloseReq(t.ctls, t.id)
  t.ctls.Close()
  return err
}

func newTransport(c *Client, ctls xnet.Stream, t *pb.Transport) (*transport, error) {
  if !t.Valid() {
    return nil, xrpc.ErrInvalidResponse
  }
  return &transport{
    id:     *t.Id,
    ctls:   ctls,
    client: c,
    code:   *t.Transport,
  }
}
