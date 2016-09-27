// Package xtpctlnet includes a set of interfaces for xtp clients and servers to use.
// These interfaces capture the types and how they work together. They are heavily
// based on Go's net package, and are almost the same as multiaddr-net's interfaces
// (we add Transport and Stream).
package xtpctlnet

import (
  "io"

  ma "github.com/multiformats/go-multiaddr"
  // manet "github.com/multiformats/go-multiaddr-net"
)

type Transport interface {
  // Code returns the transport string code (multiaddr)
  Code() string

  // Dial dials raddr, and returns a Conn if successful.
  Dial(raddr ma.Multiaddr) (Conn, error)

  // Dialer creates a dialer with laddr, and returns a Dialer if successful.
  Dialer(laddr ma.Multiaddr) (Dialer, error)

  // Listen attempts to listen on laddr, and returns a Listener if successful.
  Listen(laddr ma.Multiaddr) (Listener, error)

  // Close shuts down the transport stack, if relevant. (in many, this is a noop)
  Close() error
}

type Listener interface {
  // manet.Listener

  // // NetListener returns a compatible net.Listener.
  // NetListener() net.Listener

  // Accept waits for and returns the next connection to the listener.
  // Returns a Multiaddr friendly Conn
  Accept() (Conn, error)

  // Multiaddr returns the listener's (local) Multiaddr.
  Multiaddr() ma.Multiaddr

  // // Addr returns the net.Listener's network address.
  // Addr() net.Addr

  // Close closes the listener.
  // Any blocked Accept operations will be unblocked and return errors.
  Close() error
}

type Dialer interface {
  // Dial dials the given multiaddr and sets up a connection.
  Dial(raddr ma.Multiaddr) (Conn, error)

  // Multiaddr returns the dialer's local address, if any.
  Multiaddr() ma.Multiaddr

  // Close closes the Dialer.
  Close() error
}

type Conn interface {
  // LocalMultiaddr returns the local Multiaddr associated
  // with this connection
  LocalMultiaddr() ma.Multiaddr

  // RemoteMultiaddr returns the remote Multiaddr associated
  // with this connection
  RemoteMultiaddr() ma.Multiaddr

  // Dial attempts to open a new stream across Conn to the other side.
  Dial() (Stream, error)

  // Accept accepts an incoming conn.Dial from the other side.
  Accept() (Stream, error)

  // Close closes the Dialer.
  Close() error
}

type Stream interface {
  io.Reader
  io.Writer
  io.Closer

  // Conn returns the connection this stream belongs to.
  Conn() Conn
}
