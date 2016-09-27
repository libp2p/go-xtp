package xtpserver

import (

  xnet "github.com/libp2p/go-xtp-ctl/net"
  ma "github.com/multiformats/go-multiaddr"
)

type Server struct {
  Listener  xnet.Listener
  Xports    []xnet.Transport // to initialize with
  Clients   []ServerClient
}

func NewServer(addr ma.Multiaddr, xports []xnet.Transport) (*Server, error) {
  l, err := xnet.Listen(addr)
  if err != nil {
    return nil, err
  }

  return &Server{Listener: l, Xports: xports}, nil
}

func (s *Server) Close() error {
  s.Listener.Close()
  for _, c := range s.Clients {
    c.Close()
  }
  for _, t := range s.Xports {
    t.Close()
  }
  return nil
}
