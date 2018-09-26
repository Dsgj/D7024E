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

/*
*	PING (Message, error, timeout(bool))
 */
func (k *Kademlia) PING(c Contact) (*pb.Message, error, bool) {
	reqID := k.newRequest()
	key := c.ID.String()
	receiver := ContactToPeer(c)
	sender := ContactToPeer(k.rt.me)
	msg := k.netw.msgFct.NewPingMessage(reqID, key, sender, receiver, false)
	msgHandler, err := k.netw.SendMessage(&c, msg, true)
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan *pb.Message)
	msgHandler.awaitMessage(ch)

	respMsg := <-ch
	//handle timeouts
	/*
	*	Note updating of buckets needs to happen outside this function since
	*	PING is used when updating buckets to check if nodes are alive
	 */
	return respMsg, nil, false
}

/*
 * Note: the name is misleading; key is typically the id of the recipient
 * and we are asking for the k closest nodes that the recipient knows to be
 * closest to the key
 * a more fitting name would be FIND_CLOSE_NODES
 */
func (k *Kademlia) FIND_NODE(recipient Contact, key string) ([]*Contact, error) {
	reqID := k.newRequest()
	receiver := ContactToPeer(recipient)
	sender := ContactToPeer(k.rt.me)
	msg := k.netw.msgFct.NewFindNodeMessage(reqID, key, sender, receiver, false)
	msgHandler, err := k.netw.SendMessage(&recipient, msg, true)
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan *pb.Message)
	msgHandler.awaitMessage(ch)

	respMsg := <-ch
	//handle timeouts
	closestContacts := PeersToContacts(respMsg.GetData().GetClosestPeers())

	return closestContacts, nil
}

func (k *Kademlia) Update(c *Contact) {
	bucketIndex := k.rt.getBucketIndex(c.ID)
	bucket := k.rt.buckets[bucketIndex]
	bucket.mutex.Lock()
	if bucket.Len() < 20 { // k = 20
		bucket.AddContact(c)
		bucket.mutex.Unlock()
	} else { // bucket is full
		//ping head of bucket
		head := bucket.list.Front()
		bucket.mutex.Unlock()
		_, err, timeout := k.PING(head.Value.(Contact))
		if err != nil {
			log.Fatal(err)
		}
		if timeout {
			bucket.mutex.Lock()
			bucket.list.Remove(head)
			bucket.list.PushBack(c)
			bucket.mutex.Unlock()
		}
		//if it responds, do nothing
	}
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
