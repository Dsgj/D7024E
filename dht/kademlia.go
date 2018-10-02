package dht

import (
	pb "D7024E/dht/pb"
	"log"
	"time"
)

type Kademlia struct {
	rt        *RoutingTable
	netw      *Network
	reqCount  int32
	scheduler *Scheduler
	dataStore *Store
}

func NewKademlia(me Contact, port string) *Kademlia {
	rt := NewRoutingTable(me)
	// hardcoded port for now
	netw := NewNetwork(port, me.Address)
	k := &Kademlia{rt: rt,
		netw:      netw,
		reqCount:  0,
		scheduler: &Scheduler{},
		dataStore: NewStore()}
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
	respCh := make(chan *pb.Message)
	timeoutCh, err := k.netw.SendRequest(&c, msg, respCh)
	if err != nil {
		return nil, err, false
	}
	select {
	case respMsg := <-respCh:
		return respMsg, nil, false
	case <-time.After(30 * time.Second):
		timeoutCh <- reqID
		return nil, nil, true
	}
	/*
	*	Note updating of buckets needs to happen outside this function since
	*	PING is used when updating buckets to check if nodes are alive
	 */
}

/*
 * Note: the name is misleading; key is typically the id of the recipient
 * and we are asking for the k closest nodes that the recipient knows to be
 * closest to the key
 * a more fitting name would be FIND_CLOSE_NODES
 */
func (k *Kademlia) FIND_NODE(recipient Contact, key string) ([]Contact, error, bool) {
	reqID := k.newRequest()
	receiver := ContactToPeer(recipient)
	sender := ContactToPeer(k.rt.me)
	msg := k.netw.msgFct.NewFindNodeMessage(reqID, key, sender, receiver, false)
	respCh := make(chan *pb.Message)
	timeoutCh, err := k.netw.SendRequest(&recipient, msg, respCh)
	if err != nil {
		return nil, err, false
	}

	select {
	case respMsg := <-respCh:
		closestContacts := PeersToContacts(respMsg.GetData().GetClosestPeers())
		return closestContacts, nil, false
	case <-time.After(30 * time.Second):
		timeoutCh <- reqID
		return nil, nil, true
	}
}

func (k *Kademlia) STORE(c Contact, rec *Record, publish bool) error {
	reqID := k.newRequest()
	receiver := ContactToPeer(c)
	sender := ContactToPeer(k.rt.me)
	key := c.ID.String()
	msg := k.netw.msgFct.NewStoreMessage(reqID, key, sender, receiver, false)
	msg.AddRecord(k.dataStore.SendableRecord(GetKey(rec.value), publish))
	err := k.netw.SendMessage(&c, msg)
	if err != nil {
		return err
	}
	return nil
}

func (k *Kademlia) FINDVALUE(recipient Contact,
	key string) ([]byte, []Contact, error, bool) {
	reqID := k.newRequest()
	receiver := ContactToPeer(recipient)
	sender := ContactToPeer(k.rt.me)
	msg := k.netw.msgFct.NewFindNodeMessage(reqID, key, sender, receiver, false)
	respCh := make(chan *pb.Message)
	timeoutCh, err := k.netw.SendRequest(&recipient, msg, respCh)
	if err != nil {
		return nil, nil, err, false
	}

	select {
	case respMsg := <-respCh:
		rec := respMsg.GetData().GetRecord()
		if rec != nil {
			//we got a record
			return rec.GetValue(), nil, nil, false
		}
		// no record, take closest contacts
		closestContacts := PeersToContacts(respMsg.GetData().GetClosestPeers())
		return nil, closestContacts, nil, false
	case <-time.After(30 * time.Second):
		timeoutCh <- reqID
		return nil, nil, nil, true
	}
}

func (k *Kademlia) Update(c Contact) {
	if c.ID.String() == k.rt.me.ID.String() { //dont add yourself
		log.Printf("Attemped to add self to bucket")
		return
	}
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

func (k *Kademlia) StartScheduler() {
	task := func() {
		log.Printf("current contacts: %v", k.rt.FindClosestContacts(k.rt.me.ID, 20, k.rt.me))
		// refresh buckets
		// replicate each key/value every 1h
		// republish key/value if original publisher every 24h
		// expire key/value after 24h if not pinned
	}
	go k.scheduler.RepeatTask(10, task)
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
