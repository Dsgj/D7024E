package dht

import (
	pb "D7024E/dht/pb"
)

type rpcHandler func(*pb.Message) (*pb.Message, error)

func (k *Kademlia) getTypeHandler(t pb.Message_MessageType) rpcHandler {
	switch t {
	case pb.Message_PING:
		return k.handlePING
	case pb.Message_FIND_NODE:
		return k.handleFINDNODE
	case pb.Message_STORE:
		return k.handleSTORE
	default:
		return nil
	}
}

func (k *Kademlia) handlePING(msg *pb.Message) (*pb.Message, error) {
	/*
	 *   NewPingMessage(reqID, key, sender, receiver, isResponse)
	 *   Note flipped receiver and sender when creating response msg
	 */
	respMsg := k.netw.msgFct.NewPingMessage(msg.GetRequestID(), msg.GetKey(),
		msg.GetReceiver(), msg.GetSender(), true)
	return respMsg, nil
}

func (k *Kademlia) handleFINDNODE(msg *pb.Message) (*pb.Message, error) {
	//get nodes from bucket
	// k = 20
	target := NewKademliaID(msg.GetKey())
	// Note: Need to make sure this doesnt include the requester
	// 		(a node should never add itself to its routingtable)
	//		the requester should also discard the contact in case it happens
	//		anyway
	contacts := k.rt.FindClosestContacts(target, 20, PeerToContact(msg.GetSender()))
	peers := ContactsToPeers(contacts)
	respMsg := k.netw.msgFct.NewFindNodeMessage(msg.GetRequestID(), msg.GetKey(),
		msg.GetReceiver(), msg.GetSender(), true)
	respMsg.AddPeerData(peers)
	return respMsg, nil
}

func (k *Kademlia) handleSTORE(msg *pb.Message) (*pb.Message, error) {
	data := msg.GetData().GetRecord().GetValue()
	_, exists := k.dataStore.GetRecord(Hash(data))
	if exists {
		// what do
	} else {
		k.dataStore.Store(data)
	}
	return nil, nil
}
