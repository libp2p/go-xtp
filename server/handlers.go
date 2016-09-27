package xtpctlserver

import (
  pb "github.com/libp2p/go-xtp-ctl/pb"

  proto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
)

func (sc *ServerClient) handleReq(s IoStream, req *pb.RPC) error {
  if *(req.Error) != "" { // should not have an error in a request.
    return ErrProtocol
  }

  switch *req.Rpc {
  case pb.RPC_NoOp:
    return nil // do nothing
  case pb.RPC_ListReq:
    req2 := pb.ListReq{}
    if err := proto.Unmarshal(req.Message, &req2); err != nil {
      return err
    }
    return sc.handleListReq(s, &req2)
  case pb.RPC_CloseReq:
    req2 := pb.CloseReq{}
    if err := proto.Unmarshal(req.Message, &req2); err != nil {
      return err
    }
    return sc.handleCloseReq(s, &req2)
  case pb.RPC_ListenReq:
    return sc.handleListenReq(s, req)
  case pb.RPC_AcceptReq:
    return sc.handleAcceptReq(s, req)
  case pb.RPC_DialerReq:
    return sc.handleDialerReq(s, req)
  case pb.RPC_DialReq:
    return sc.handleDialReq(s, req)
  default:
    return ErrUnknownRPC
  }
}

func (sc *ServerClient) handleListReq(s IoStream, req *pb.ListReq) error {
  res := pb.ListRes{}

  types := req.TypesRequested()

  addItem := func(i *pb.ListRes_Item, err error) {
    if err != nil {
      res.Items = append(res.Items, i)
    } else {
      // internal error. TODO: log it
    }
  }

  sc.lk.Lock()
  for _, t := range sc.Transports {
    t.Lock()

    // add the transport to the list of items.
    if types.Transports {
      i, err := pb.ListRes_Item_Transport(t.PB())
      addItem(i, err)
    }

    // add the transport listeners
    if types.Listeners {
      for _, l := range t.Listeners {
        i, err := pb.ListRes_Item_Listener(l.PB())
        addItem(i, err)
      }
    }

    // add the transport dialers
    if types.Dialers {
      for _, d := range t.Dialers {
        i, err := pb.ListRes_Item_Dialer(d.PB())
        addItem(i, err)
      }
    }

    if types.Conns || types.Streams {
      for _, c := range t.Conns {
        c.Lock()

        // add the conn
        if types.Conns {
          i, err := pb.ListRes_Item_Conn(c.PB())
          addItem(i, err)
        }

        // add the streams
        if types.Streams {
          for _, s := range c.Streams {
            i, err := pb.ListRes_Item_Stream(s.PB())
            addItem(i, err)
          }
        }

        c.Unlock()
      }
    }

    t.Unlock()
  }

  sc.lk.Unlock() // todo: more granular locking, to avoid holding lock while marshalling.

  return StreamWriteRPCRes(s, pb.RPC_ListRes, &res)
}

func (sc *ServerClient) handleCloseReq(s IoStream, req *pb.CloseReq) error {
  id := *req.Id
  sc.CloseId(id)
  return StreamWriteRPCRes(s, pb.RPC_CloseRes, nil)
}

func (sc *ServerClient) handleListenReq(s IoStream, req *pb.RPC) error {
  panic("not implemented yet")
}

func (sc *ServerClient) handleAcceptReq(s IoStream, req *pb.RPC) error {
  panic("not implemented yet")
}

func (sc *ServerClient) handleDialerReq(s IoStream, req *pb.RPC) error {
  panic("not implemented yet")
}

func (sc *ServerClient) handleDialReq(s IoStream, req *pb.RPC) error {
  panic("not implemented yet")
}