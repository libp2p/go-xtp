package xtpclient

import (
  ma "github.com/multiformats/go-multiaddr"
  manet "github.com/multiformats/go-multiaddr-net"
  proto "github.com/gogo/protobuf/proto"
  xnet "github.com/libp2p/go-xtp-ctl/net"
  xrpc "github.com/libp2p/go-xtp-ctl/rpc"
  pb "github.com/libp2p/go-xtp-ctl/pb"
)

type Client struct {
  Conn   xnet.Conn // connection to the server
  Xports []xnet.Transport
}

func NewClient(server ma.Multiaddr) (*Client, error) {
  c, err := manet.Dial(server)
  if err != nil {
    return nil, err
  }

  // layer muxer
  c2, err := xnet.XtpCtlConn(c, false)
  if err != nil {
    return nil, err
  }

  client := &Client{Conn: c2}
  // first, figure out the transports
  if err := client.getTransports(); err != nil {
    client.Close()
    return nil, err
  }
  return client, nil
}

func (c *Client) getTransports() error {
  s, err := c.Conn.Dial()
  if err != nil {
    return err
  }

  items, err := xrpc.ListReq(s, []pb.TType{pb.TType_TTypeTransport})
  if err != nil {
    return err
  }

  var xports []xnet.Transport
  for _, item := range items {
    if *item.Type != pb.TType_TTypeTransport {
      continue // skip
    }

    pbt := &pb.Transport{}
    err := proto.Unmarshal(item.Value, pbt)
    if err != nil {
      continue // skip
    }

    t, err := newTransport(c, nil, pbt)
    if err != nil {
      continue // skip
    }
    xports = append(xports, t)
  }
  c.Xports = xports
  return nil
}

func (c *Client) Close() error {
  return c.Conn.Close()
}
