package xtpimpls

import (
  "errors"
  "sync"

  ma "github.com/multiformats/go-multiaddr"
  manet "github.com/multiformats/go-multiaddr-net"
  xnet "github.com/libp2p/go-xtp-ctl/net"
)

var ErrNoMoreStreams = errors.New("unable to create more streams")

type transport struct {
  code string
}


func (t *transport) Code() string { return t.code }
func (t *transport) Dial(raddr ma.Multiaddr) (xnet.Conn, error) {
  c, err := manet.Dial(raddr)
  return &conn{C: c}, err
}

func (t *transport) Dialer(laddr ma.Multiaddr) (xnet.Dialer, error) {
  d := manet.Dialer{LocalAddr: laddr}
  return &dialer{d}, nil
}

func (t *transport) Listen(laddr ma.Multiaddr) (xnet.Listener, error) {
  l, err := manet.Listen(laddr)
  return &listener{l}, err
}

func (t *transport) Close() error {
  return nil
}


type listener struct {
  L manet.Listener
}

func (l *listener) Accept() (xnet.Conn, error) {
  c, err := l.L.Accept()
  return &conn{C: c}, err
}

func (l *listener) Multiaddr() ma.Multiaddr { return l.L.Multiaddr() }
func (l *listener) Close() error { return l.L.Close() }

type dialer struct {
  D manet.Dialer
}

func (d *dialer) Dial(raddr ma.Multiaddr) (xnet.Conn, error) {
  c, err := d.D.Dial(raddr)
  return &conn{C: c}, err
}
func (d *dialer) Multiaddr() ma.Multiaddr { return d.D.LocalAddr }
func (d *dialer) Close() error { return nil }

type conn struct {
  C manet.Conn

  lk sync.Mutex
  d  bool
  l  bool
}

func (c *conn) LocalMultiaddr() ma.Multiaddr { return c.C.LocalMultiaddr() }
func (c *conn) RemoteMultiaddr() ma.Multiaddr { return c.C.RemoteMultiaddr() }
func (c *conn) Close() error { return c.C.Close() }

func (c *conn) Dial() (xnet.Stream, error) {
  c.lk.Lock()
  defer c.lk.Unlock()

  if c.d {
    return nil, ErrNoMoreStreams
  }
  c.d = true
  return &singleStream{c}, nil
}

func (c *conn) Accept() (xnet.Stream, error) {
  c.lk.Lock()
  defer c.lk.Unlock()

  if c.l {
    return nil, ErrNoMoreStreams
  }
  c.l = true
  return &singleStream{c}, nil
}

type singleStream struct {
  C *conn
}

func (s *singleStream) Conn() xnet.Conn {
  return s.C
}

func (s *singleStream) Read(buf []byte) (int, error) {
  return s.C.Read(buf)
}

func (s *singleStream) Write(buf []byte) (int, error) {
  return s.C.Write(buf)
}

func (s *singleStream) Close() error {
  return s.C.Close()
}

