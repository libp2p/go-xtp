package xtpctlclient

import (
  "io"

  pb "github.com/libp2p/go-xtp-ctl/pb"
  xnet "github.com/libp2p/go-xtp-ctl/net"
  xrpc "github.com/libp2p/go-xtp-ctl/rpc"
)

type stream struct {
  id     int64
  tid    int64
  ctls   IoStream // the xtp-ctl stream for this conn
  client *Client  // the xtp-ctl client
  conn   *conn    // the conn this stream belongs to
}

func (s *stream) Read(buf []byte) (int, error) {
  return s.ctls.Read(buf)
}

func (s *stream) Write(buf []byte) (int, error) {
  return s.ctls.Write(buf)
}

// Conn returns the Conn this stream belongs to.
func (s *stream) Conn() xnet.Conn {
  return s.conn
}


// Close closes the dialer.
func (s *stream) Close() error {
  err := xrpc.CloseReq(s.ctls, s.id)
  s.ctls.Close()
  return err
}

func newStream(c *Client, ctls xnet.Stream, s *pb.Stream, cn *conn) (*stream, error) {
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
