package xtpserver

import (
  "io"
  "errors"

  pb "github.com/libp2p/go-xtp-ctl/pb"
  ma "github.com/multiformats/go-multiaddr"
  xrpc "github.com/libp2p/go-xtp-ctl/rpc"

  proto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
)

func rpcHandler(sc *ServerClient, s IoStream) error {
  req := &pb.RPC{}
  if err := xrpc.ReadRPC(s, req); err != nil {
    return err
  }

  err := handleReq(sc, s, req)
  if err != nil {
    return xrpc.ErrRPCRes(s, req, err)
  }
  return nil
}

func handleReq(sc *ServerClient, s IoStream, req *pb.RPC) error {
  if *(req.Error) != "" { // should not have an error in a request.
    return xrpc.ErrProtocol
  }

  switch *req.Rpc {
  case pb.RPC_NoOp:
    return nil // do nothing
  case pb.RPC_ListReq:
    req2 := &pb.ListReq{}
    if err := proto.Unmarshal(req.Message, req2); err != nil {
      return err
    }
    return handleListReq(sc, s, req2)
  case pb.RPC_CloseReq:
    req2 := &pb.CloseReq{}
    if err := proto.Unmarshal(req.Message, req2); err != nil {
      return err
    }
    return handleCloseReq(sc, s, req2)
  case pb.RPC_ListenReq:
    req2 := &pb.ListenReq{}
    if err := proto.Unmarshal(req.Message, req2); err != nil {
      return err
    }
    return handleListenReq(sc, s, req2)
  case pb.RPC_AcceptReq:
    req2 := &pb.AcceptReq{}
    if err := proto.Unmarshal(req.Message, req2); err != nil {
      return err
    }
    return handleAcceptReq(sc, s, req2)
  case pb.RPC_DialerReq:
    req2 := &pb.DialerReq{}
    if err := proto.Unmarshal(req.Message, req2); err != nil {
      return err
    }
    return handleDialerReq(sc, s, req2)
  case pb.RPC_DialReq:
    req2 := &pb.DialReq{}
    if err := proto.Unmarshal(req.Message, req2); err != nil {
      return err
    }
    return handleDialReq(sc, s, req2)
  default:
    return xrpc.ErrUnknownRPC
  }
}

func handleListReq(sc *ServerClient, s IoStream, req *pb.ListReq) error {
  types := req.TypesRequested()

  var items []*pb.ListRes_Item
  addItem := func(i *pb.ListRes_Item, err error) {
    if err != nil {
      items = append(items, i)
    } else {
      // internal error. TODO: log it
    }
  }

  sc.Lock()
  for _, t := range sc.transports {
    // add the transport to the list of items.
    if types.Transports {
      i, err := pb.ListRes_Item_Transport(t.PB())
      addItem(i, err)
    }

    if types.Listeners || types.Dialers || types.Conns || types.Streams {
      items = append(items, t.List(types)...)
    }
  }
  sc.Unlock() // todo: more granular locking, to avoid holding lock while marshalling.

  res := pb.ListRes{Items: items}
  return xrpc.WriteRPCMsg(s, pb.RPC_ListRes, &res, nil)
}

func handleCloseReq(sc *ServerClient, s IoStream, req *pb.CloseReq) error {
  id := *req.Id
  if id < pb.MinId {
    return xrpc.ErrInvalidMessage
  }

  err := sc.CloseId(id)
  if err != nil {
    // log error, internal failure.
  }

  return xrpc.WriteRPCMsg(s, pb.RPC_CloseRes, nil, nil)
}

func handleListenReq(sc *ServerClient, s IoStream, req *pb.ListenReq) error {
  l := req.ListenerOpts
  if l == nil || l.Multiaddr == nil || l.TransportId == nil {
    return xrpc.ErrInvalidMessage
  }

  // get parameters
  laddr, err := ma.NewMultiaddrBytes(l.Multiaddr)
  if err != nil {
    return err
  }

  tid := *l.TransportId
  t := sc.transport(tid)
  if t == nil {
    return errors.New("transport id not found")
  }

  // listen
  l2, err := t.Listen(laddr)
  if err != nil {
    return err
  }

  // send response with listener
  return xrpc.ListenRes(s, l2.PB(), nil)
}

func handleAcceptReq(sc *ServerClient, s IoStream, req *pb.AcceptReq) error {
  if req.Id == nil {
    return xrpc.ErrInvalidMessage
  }

  // get parameters
  id := *req.Id

  v := sc.Find(id)

  var c1 *pb.Conn
  var s1 *pb.Stream

  switch v := v.(type) {
  case *listener:
    c2, err := v.Accept()
    if err != nil {
      return err
    }
    c1 = c2.PB()
  case *conn:
    s2, err := v.Accept()
    if err != nil {
      return err
    }
    s1 = s2.PB()
  default:
    return errors.New("id mismatch (not a listener or conn)")
  }

  // send response with listener
  return xrpc.AcceptRes(s, c1, s1, nil)
}

func handleDialerReq(sc *ServerClient, s IoStream, req *pb.DialerReq) error {
  d := req.DialerOpts
  if d == nil || d.Multiaddr == nil || d.TransportId == nil {
    return xrpc.ErrInvalidMessage
  }

  // get parameters
  laddr, err := ma.NewMultiaddrBytes(d.Multiaddr)
  if err != nil {
    return err
  }

  tid := *d.TransportId
  t := sc.transport(tid)
  if t == nil {
    return errors.New("transport id not found")
  }

  // dial
  d2, err := t.Dialer(laddr)
  if err != nil {
    return err
  }

  // send response with dialer
  return xrpc.DialerRes(s, d2.PB(), nil)
}

func handleDialReq(sc *ServerClient, s IoStream, req *pb.DialReq) error {
 if req.Id == nil {
    return xrpc.ErrInvalidMessage
  }

  // get parameters
  id := *req.Id
  opts := *req.ConnOpts

  v := sc.Find(id)

  var c1 *pb.Conn
  var s1 *pb.Stream

  switch v := v.(type) {
  case *dialer:
    raddr, err := ma.NewMultiaddrBytes(opts.LocalMultiaddr)
    if err != nil {
      return err
    }
    c2, err := v.Dial(raddr)
    if err != nil {
      return err
    }
    c1 = c2.PB()
  case *conn:
    s2, err := v.Dial()
    if err != nil {
      return err
    }
    s1 = s2.PB()
  default:
    return errors.New("id mismatch (not a listener or conn)")
  }

  // send response with listener
  return xrpc.DialRes(s, c1, s1, nil)
}

type IoStream interface {
  io.Reader
  io.Writer
  io.Closer
}