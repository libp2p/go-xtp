package xtpctlrpc

import (
  "errors"
  "io"
  "fmt"

  ggio "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/io"
  proto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"

  pb "github.com/libp2p/go-xtp-ctl/pb"
)

var (
  ErrUnknownRPC = errors.New("unknown rpc")
  ErrProtocol   = errors.New("incorrect protocol behavior")
  ErrNotFound   = errors.New("descriptor not found")
  ErrInvalidResponse = errors.New("invalid response")
)

var (
  MessageSizeMax = 1 << 12
)

func WriteRPC(s IoStream, rpc *pb.RPC) error {
  w := ggio.NewDelimitedWriter(s)
  return w.WriteMsg(rpc)
}

func ReadRPC(s IoStream, rpc *pb.RPC) error {
  r := ggio.NewDelimitedReader(s, MessageSizeMax)
  return r.ReadMsg(rpc)
}

func WriteRPCMsg(s IoStream, typ pb.RPC_Type, m proto.Message, err error) error {
  rpc := pb.RPC{Rpc: &typ}
  if err != nil {
    estr := fmt.Sprint(err)
    rpc.Error = &estr
  }
  if m != nil {
    b, err := proto.Marshal(m)
    if err != nil {
      return err
    }
    rpc.Message = b
  }
  return WriteRPC(s, &rpc)
}

func ReadRPCMsg(s IoStream, typ pb.RPC_Type, m proto.Message) error {
  rpc := pb.RPC{}
  err := ReadRPC(s, &rpc)
  if err != nil {
    return err
  }

  if !rpc.Valid() {
    return ErrInvalidResponse
  }
  if typ != pb.RPC_Null && typ != *rpc.Rpc {
    return ErrInvalidResponse
  }

  if rpc.Error != nil && len(*rpc.Error) > 0 {
    return errors.New(*rpc.Error)
  }

  // ok.
  if m == nil {
    return nil // not interested in return message
  }
  if rpc.Message == nil || len(rpc.Message) < 1 {
    return nil // no return message
  }
  if err := proto.Unmarshal(rpc.Message, m); err != nil {
    return err
  }

  // automatic validation :)
  if v, ok := m.(validator); ok {
    if !v.Valid() {
      return ErrInvalidResponse
    }
  }
  return nil
}

type IoStream interface {
  io.Reader
  io.Writer
  io.Closer
}

type validator interface {
  Valid() bool
}
