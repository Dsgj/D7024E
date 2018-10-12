package dht

import (
	pb "D7024E/dht/pb"
	"log"
	"time"
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
	case pb.Message_FIND_VALUE:
		return k.handleFINDVALUE
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
	key := GetKey(data)
	rec, exists := k.dataStore.GetRecord(key)
	// testing
	log.Printf("publisher(peer): %+v\n", msg.GetData().GetRecord().GetPublisher())
	log.Printf("publish: %+v\n", msg.GetData().GetRecord().GetNewPublish())
	log.Printf("value: %+v\n", msg.GetData().GetRecord().GetValue())
	publisher := PeerToContact(msg.GetData().GetRecord().GetPublisher())
	publAt := time.Unix(msg.GetData().GetRecord().GetPublishedAt(), 0)
	if exists { // record exists in local store
		if msg.GetData().GetRecord().GetNewPublish() { //republish, update time
			rec.Republish(publAt)
		}
	} else { // new record
		k.dataStore.Store(data, publisher, publAt)
	}
	return nil, nil
}

func (k *Kademlia) handleFINDVALUE(msg *pb.Message) (*pb.Message, error) {
	key, err := FromString(msg.GetKey())
	if err != nil {
		return nil, err
	}
	respMsg := k.netw.msgFct.NewFindValueMessage(msg.GetRequestID(), msg.GetKey(),
		msg.GetReceiver(), msg.GetSender(), true)

	_, exists := k.dataStore.GetRecord(key)
	if exists {
		respMsg.AddRecord(k.dataStore.SendableRecord(key, false))
	} else {
		target := NewKademliaID(msg.GetKey())
		contacts := k.rt.FindClosestContacts(target, 20, PeerToContact(msg.GetSender()))
		peers := ContactsToPeers(contacts)
		respMsg.AddPeerData(peers)
	}

	return respMsg, nil
}
