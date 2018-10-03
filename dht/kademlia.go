package dht

import (
	pb "D7024E/dht/pb"
	"log"
	"math/rand"
	"time"
)

const requestTimeout = 5 * time.Second
const alpha = 3

type Kademlia struct {
	rt        *RoutingTable
	netw      *Network
	reqCount  int32
	scheduler *Scheduler
	dataStore *Store
}

func NewKademlia(me Contact, port string) *Kademlia {
	rt := NewRoutingTable(me)
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
	case <-time.After(requestTimeout):
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
	case <-time.After(requestTimeout):
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

func (k *Kademlia) FIND_VALUE(recipient Contact,
	key string) ([]byte, []Contact, error, bool) {
	reqID := k.newRequest()
	receiver := ContactToPeer(recipient)
	sender := ContactToPeer(k.rt.me)
	msg := k.netw.msgFct.NewFindValueMessage(reqID, key, sender, receiver, false)
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
	case <-time.After(requestTimeout):
		timeoutCh <- reqID
		return nil, nil, nil, true
	}
}

func (k *Kademlia) IterativeStore(key [20]byte, publish bool) {
	closestContacts, err := k.IterativeFindNode(ToString(key))
	if err != nil {
		log.Println(err)
	}
	for _, contact := range closestContacts {
		reqID := k.newRequest()
		receiver := ContactToPeer(contact)
		sender := ContactToPeer(k.rt.me)
		msg := k.netw.msgFct.NewStoreMessage(reqID, ToString(key), sender, receiver, false)
		msg.AddRecord(k.dataStore.SendableRecord(key, publish))
		err := k.netw.SendMessage(&contact, msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func (k *Kademlia) IterativeFindNode(key string) ([]Contact, error) {
	toBeQueried := k.rt.FindClosestContacts(NewKademliaID(key), 20, k.rt.me)
	log.Printf("to be queried: %+v\n", toBeQueried)
	alreadyQueried := make(map[*KademliaID]bool)
	shortList := &ContactCandidates{}

	for {
		countNodesToQuery := 0
		alreadyAdded := make(map[*KademliaID]bool)
		for i := 0; i < alpha; i++ {
			for _, contact := range toBeQueried {
				if !alreadyQueried[contact.ID] && !alreadyAdded[contact.ID] {
					if !shortList.Exists(contact) {
						shortList.Add(contact)
						alreadyAdded[contact.ID] = true
						countNodesToQuery++
					}
				}
			}
		}
		log.Printf("current shortlist: %+v\n", shortList)
		if countNodesToQuery == 0 { // we queried all nodes
			shortList.Sort()
			shortList.Cut()
			return shortList.contacts, nil
		} else {
			_, shortList, alreadyQueried = k.findCloserNodesOrValue(shortList, key, alreadyQueried, false)

		}

	}
	return shortList.contacts, nil
}

func (k *Kademlia) IterativeFindValue(key string) ([]byte, []Contact, error) {
	toBeQueried := k.rt.FindClosestContacts(NewKademliaID(key), 20, k.rt.me)
	log.Printf("to be queried: %+v\n", toBeQueried)
	alreadyQueried := make(map[*KademliaID]bool)
	shortList := &ContactCandidates{}
	var value []byte

	for {
		countNodesToQuery := 0
		alreadyAdded := make(map[*KademliaID]bool)
		for i := 0; i < alpha; i++ {
			for _, contact := range toBeQueried {
				if !alreadyQueried[contact.ID] && !alreadyAdded[contact.ID] {
					if !shortList.Exists(contact) {
						shortList.Add(contact)
						alreadyAdded[contact.ID] = true
						countNodesToQuery++
					}
				}
			}
		}
		log.Printf("current shortlist: %+v\n", shortList)
		if countNodesToQuery == 0 { // we queried all nodes
			shortList.Sort()
			shortList.Cut()
			return nil, shortList.contacts, nil
		} else {
			value, shortList, alreadyQueried = k.findCloserNodesOrValue(shortList, key, alreadyQueried, true)
			if value != nil {
				return value, shortList.contacts, nil
			}
		}

	}
	return nil, shortList.contacts, nil
}

func (k *Kademlia) findCloserNodesOrValue(shortList *ContactCandidates,
	key string,
	alreadyQueried map[*KademliaID]bool,
	wantValue bool) ([]byte, *ContactCandidates, map[*KademliaID]bool) {
	done := make(chan []Contact)
	timeoutCh := make(chan Contact)
	valueCh := make(chan []byte)
	closestContact := shortList.contacts[0]
	pending := 0
	countNoCloserNodes := 0
	log.Printf("finding closer nodes, shortlist: %+v\n", shortList)
	for {
		select {
		case value := <-valueCh:
			return value, shortList, alreadyQueried
		case newContacts := <-done:
			shortList.AddUnique(newContacts)
			shortList.Sort()
			shortList.Cut()
			newClosestContact := shortList.contacts[0]
			if newClosestContact.Equals(&closestContact) {
				countNoCloserNodes++
			} else {
				closestContact = newClosestContact
				countNoCloserNodes = 0
			}
			if (countNoCloserNodes) >= alpha {
				return nil, shortList, alreadyQueried
			}
			pending--
		case badContact := <-timeoutCh:
			for i, contact := range shortList.contacts {
				if contact.Equals(&badContact) {
					shortList.Remove(i)
				}
			}
			pending--
		default:
			if pending < alpha {
				for _, contact := range shortList.contacts {
					if !alreadyQueried[contact.ID] {
						alreadyQueried[contact.ID] = true
						pending++
						go func() {
							var value []byte
							var contacts []Contact
							var err error
							var timeout bool
							if wantValue {
								value, contacts, err, timeout = k.FIND_VALUE(contact, key)
							} else {
								contacts, err, timeout = k.FIND_NODE(contact, key)
							}
							if err != nil {
								log.Println(err)
								return
							}
							if value != nil && wantValue {
								valueCh <- value
								return
							}
							if timeout {
								timeoutCh <- contact
								return
							}
							done <- contacts
							return
						}()
						break
					}
				}
				if pending == 0 {
					return nil, shortList, alreadyQueried
				}
			}
		}
	}
}

func (k *Kademlia) Update(c Contact) {
	if c.Equals(&k.rt.me) { //dont add yourself
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
			log.Println(err)
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
		//log.Printf("current contacts: %v", k.rt.FindClosestContacts(k.rt.me.ID, 20, k.rt.me))
		// refresh buckets
		for _, bucket := range k.rt.buckets {
			if bucket.NeedsRefresh(time.Now()) {
				contact := bucket.GetRandomContact()
				if contact != nil {
					bucket.Refresh(time.Now())
					log.Printf("refreshing bucket: %v", bucket)
					go func() {
						_, err := k.IterativeFindNode(contact.ID.String())
						if err != nil {
							log.Println(err)
						} else {
							log.Printf("bucket refreshed: %v", bucket)
						}
					}()
				}
			}
		}

		for key, record := range k.dataStore.records {
			if record.IsExpired(time.Now()) {
				k.dataStore.DelRecord(key)
			} else if record.NeedsRepublish(time.Now(), k.rt.me) {
				record.Republish(time.Now())
				go k.IterativeStore(key, true)
			} else if record.NeedsReplicate(time.Now()) {
				record.Replicate(time.Now())
				go k.IterativeStore(key, false)
			}
		}
	}
	go k.scheduler.RepeatTask(10, task)
}

func (k *Kademlia) TestStore() { //manual testing
	rand.Seed(time.Now().UnixNano())
	N := rand.Intn(10)
	testBytes := make([]byte, N)
	for i := 0; i < N; i++ {
		testBytes[i] = 'a' + byte(i%26)
	}
	rec := k.dataStore.Store(testBytes, k.rt.me, time.Now())
	log.Printf("iterativestore on rec: %v", rec)
	k.IterativeStore(GetKey(testBytes), true)
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
