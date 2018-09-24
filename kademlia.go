package d7024e

import (
	pb "D7024E/pb"
	"log"
)

type Kademlia struct {
	rt       *RoutingTable
	netw     *Network
	reqCount int32
}

func NewKademlia(me Contact, port string) *Kademlia {
	rt := NewRoutingTable(me)
	// hardcoded port for now
	netw := NewNetwork(port, me.Address)
	k := &Kademlia{rt: rt,
		netw:     netw,
		reqCount: 0}
	return k
}

func (k *Kademlia) newRequest() int32 {
	current := k.reqCount
	k.reqCount++
	return current
}

func (k *Kademlia) PING(c Contact) (*pb.Message, error) {
	reqID := k.newRequest()
	key := c.ID.String()
	receiver := ContactToPeer(c)
	//should be kademlia.self or something
	//sender is receiver for now
	sender := ContactToPeer(c)
	msg := k.netw.msgFct.NewPingMessage(reqID, key, sender, receiver, false)
	msgHandler, err := k.netw.SendMessage(&c, msg, true)
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan *pb.Message)
	msgHandler.awaitMessage(ch)

	respMsg := <-ch
	/*
	*	Note updating of buckets needs to happen outside this function since
	*	PING is used when updating buckets to check if nodes are alive
	 */
	return respMsg, nil
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
