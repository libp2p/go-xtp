package xtpctlserver

import (
  "sync"
  "fmt"

  manet "github.com/multiformats/go-multiaddr-net"
  ggio "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/io"
  proto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"

  pb "github.com/libp2p/go-xtp-ctl/pb"
)

type ServerClient struct {
  Server *Server
  Conn   manet.Conn

  lk         sync.RWMutex
  nextId     int64
  Transports map[int64]*Transport

  NextId func() int64 // function to get the next id.
}

func StreamWriteRPC(s IoStream, rpc *pb.RPC) error {
  w := ggio.NewDelimitedWriter(s)
  return w.WriteMsg(rpc)
}

func StreamWriteRPCRes(s IoStream, typ pb.RPC_Type, res proto.Message) error {
  var b []byte
  var err error
  if res != nil {
    b, err = proto.Marshal(res)
    if err != nil {
      return err
    }
  }
  return StreamWriteRPC(s, &pb.RPC{
    Rpc:     &typ,
    Message: b,
  })
}

func StreamWriteRPCError(s IoStream, typ pb.RPC_Type, err error) error {
  str := fmt.Sprint(err)
  return StreamWriteRPC(s, &pb.RPC{
    Rpc:   &typ,
    Error: &str,
  })
}

func (sc *ServerClient) RPCStreamHandler(s IoStream) error {
  r := ggio.NewDelimitedReader(s, MessageSizeMax)
  req := pb.RPC{}
  if err := r.ReadMsg(&req); err != nil {
    return err
  }

  err := sc.handleReq(s, &req)
  if err != nil {
    return err
  }

  return nil
}

// Close shuts down the ServerClient, closing everything.
func (sc *ServerClient) Close() error {
  panic("todo")
}

func (sc *ServerClient) closeTransport(id int64) (err error) {
  sc.lk.Lock()
  defer sc.lk.Unlock()
  t, found := sc.Transports[id]
  if found {
    err = t.Close()
    delete(sc.Transports, id)
  }
  return err
}

func (sc *ServerClient) CloseId(id int64) error {
  sc.lk.RLock()
  _, found := sc.Transports[id]
  if !found {
    for _, t := range sc.Transports {
      t.CloseId(id)
    }
  }
  sc.lk.RUnlock()

  if found { // it's a transport, remove it here (different Lock, not RLock)
    sc.closeTransport(id)
  }

  return nil
}

func (sc *ServerClient) AddTransport(code string) *Transport {
  t := NewTransport(sc.NextId, code)
  sc.lk.Lock()
  sc.Transports[t.Id] = t
  sc.lk.Unlock()
  return t
}
