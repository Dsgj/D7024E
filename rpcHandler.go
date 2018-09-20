package d7024e

import (
    pb "D7024E/pb"
)

type rpcHandler func(*pb.Message) (*pb.Message, error)

func (k *Kademlia) getTypeHandler (t pb.Message_MessageType) rpcHandler {
    switch t {
    case pb.Message_PING:
        return k.handlePING
    default:
        return nil
    }
}

func (k *Kademlia) handlePING(msg *pb.Message) (*pb.Message, error) {
    //return response msg
    return nil, nil
}
