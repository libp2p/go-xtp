package xtpctlrpc

import (
  pb "github.com/libp2p/go-xtp-ctl/pb"
)

func ListReq(s IoStream, types []pb.TType) ([]*pb.ListRes_Item, error) {
  // send the request
  err := WriteRPCMsg(s, pb.RPC_ListReq, &pb.ListReq{Types: types}, nil)
  if err != nil {
    return nil, err
  }

  // now get the response
  res := pb.ListRes{}
  if err := ReadRPCMsg(s, pb.RPC_ListRes, &res); err != nil {
    return nil, err
  }
  return res.Items, nil
}

func ListRes(s IoStream, items []*pb.ListRes_Item, err error) error {
  // send the response
  return WriteRPCMsg(s, pb.RPC_ListRes, &pb.ListRes{Items: items}, err)
}

func CloseReq(s IoStream, id int64) error {
  // send the request
  err := WriteRPCMsg(s, pb.RPC_CloseReq, &pb.CloseReq{Id: &id}, nil)
  if err != nil {
    return err
  }

  // now get the response
  return ReadRPCMsg(s, pb.RPC_CloseRes, nil)
}

func ListenReq(s IoStream, lopts *pb.Listener) (listener *pb.Listener, err error) {
  // send the request
  err = WriteRPCMsg(s, pb.RPC_ListReq, &pb.ListenReq{ListenerOpts: lopts}, nil)
  if err != nil {
    return nil, err
  }

  // now get the response
  res := pb.ListenRes{}
  if err := ReadRPCMsg(s, pb.RPC_ListenRes, &res); err != nil {
    return nil, err
  }
  return res.Listener, nil
}

func ListenRes(s IoStream, l *pb.Listener, err error) error {
  return WriteRPCMsg(s, pb.RPC_ListenRes, &pb.ListenRes{Listener: l}, err)
}

func AcceptReq(s IoStream, id int64) (*pb.AcceptRes, error) {
  // send the request
  err := WriteRPCMsg(s, pb.RPC_AcceptReq, &pb.AcceptReq{Id: &id}, nil)
  if err != nil {
    return nil, err
  }

  // now get the response
  res := pb.AcceptRes{}
  if err := ReadRPCMsg(s, pb.RPC_AcceptRes, &res); err != nil {
    return nil, err
  }
  return &res, nil
}

func AcceptRes(s IoStream, conn *pb.Conn, st *pb.Stream, err error) error {
  return WriteRPCMsg(s, pb.RPC_AcceptRes, &pb.AcceptRes{Conn: conn, Stream: st}, err)
}

// todo: connOpts
func DialerReq(s IoStream, dopts *pb.Dialer) (d *pb.Dialer, err error) {
  // send the request
  err = WriteRPCMsg(s, pb.RPC_DialerReq, &pb.DialerReq{DialerOpts: dopts}, nil)
  if err != nil {
    return nil, err
  }

  // now get the response
  res := pb.DialerRes{}
  if err := ReadRPCMsg(s, pb.RPC_DialerRes, &res); err != nil {
    return nil, err
  }
  return res.Dialer, nil
}

func DialerRes(s IoStream, d *pb.Dialer, err error) error {
  return WriteRPCMsg(s, pb.RPC_DialerRes, &pb.DialerRes{Dialer: d}, err)
}

func DialReq(s IoStream, id int64) (*pb.DialRes, error) {
  // send the request
  err := WriteRPCMsg(s, pb.RPC_DialReq, &pb.DialReq{Id: &id}, nil)
  if err != nil {
    return nil, err
  }

  // now get the response
  res := pb.DialRes{}
  if err := ReadRPCMsg(s, pb.RPC_DialRes, &res); err != nil {
    return nil, err
  }
  return &res, nil
}

func DialRes(s IoStream, conn *pb.Conn, st *pb.Stream, err error) error {
  return WriteRPCMsg(s, pb.RPC_DialReq, &pb.DialRes{Conn: conn, Stream: st}, err)
}
