package xtpctlnet

import (
  ma "github.com/multiformats/go-multiaddr"
  manet "github.com/multiformats/go-multiaddr-net"
  ymux "gx/ipfs/QmSHTSkxXGQgaHWz91oZV3CDy3hmKmDgpjbYRT6niACG4E/go-smux-yamux"
  smux "gx/ipfs/Qmb1US8uyZeEpMyc56wVZy2cDFdQjNFojAUYVCoo9ieTqp/go-stream-muxer"
)

// XtpCtlConn wraps a raw manet.Conn with the necessary
// protocols for XTP-Ctl. For now this means:
// - yamux
// the server parameter is used by yamux.
func XtpCtlConn(c manet.Conn, server bool) (Conn, error) {
  tr := ymux.DefaultTransport
  sc, err := tr.NewConn(c, server)
  return &smuxConn{c, sc}, err
}

type smuxConn struct {
  C manet.Conn
  S smux.Conn
}

func (c *smuxConn) LocalMultiaddr() ma.Multiaddr {
  return c.C.LocalMultiaddr()
}

func (c *smuxConn) RemoteMultiaddr() ma.Multiaddr {
  return c.C.RemoteMultiaddr()
}

func (c *smuxConn) Dial() (Stream, error) {
  s, err := c.S.OpenStream()
  return &smuxStream{c, s}, err
}

func (c *smuxConn) Accept() (Stream, error) {
  s, err := c.S.AcceptStream()
  return &smuxStream{c, s}, err
}

func (c *smuxConn) Close() error {
  return c.S.Close()
}

type smuxStream struct {
  C Conn
  S smux.Stream
}

func (s *smuxStream) Conn() Conn {
  return s.C
}

func (s *smuxStream) Read(buf []byte) (int, error) {
  return s.S.Read(buf)
}

func (s *smuxStream) Write(buf []byte) (int, error) {
  return s.S.Write(buf)
}

func (s *smuxStream) Close() error {
  return s.S.Close()
}
