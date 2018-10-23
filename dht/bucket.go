package dht

import (
	"container/list"
	"log"
	"math/rand"
	"sync"
	"time"
)

//const tRefresh = time.Hour * 1
const tRefresh = time.Second * 20

// bucket definition
// contains a List
type bucket struct {
	list       *list.List
	mutex      *sync.Mutex
	lastUpdate time.Time
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	bucket.mutex = &sync.Mutex{}
	bucket.lastUpdate = time.Now()
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact) {
	bucket.lastUpdate = time.Now()
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
			log.Printf("added contact to bucket: %+v\n", contact)
		}
	} else {
		bucket.list.MoveToFront(element)
		//log.Printf("moved contact to front of bucket: %+v\n", contact)
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID, ignore Contact) []Contact {
	var contacts []Contact

	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()
	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		if !contact.Equals(&ignore) {
			contact.CalcDistance(target)
			contacts = append(contacts, contact)
		}
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}

func (bucket *bucket) NeedsRefresh(t time.Time) bool {
	return t.Sub(bucket.lastUpdate) >= tRefresh && bucket.Len() > 0
}

func (bucket *bucket) Refresh(t time.Time) {
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()
	bucket.lastUpdate = t
}

func (bucket *bucket) GetRandomContact() *Contact {
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()
	if bucket.Len() == 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(bucket.Len())
	j := 0
	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		if j == index {
			contact := elt.Value.(Contact)
			return &contact
		}
		j++
	}
	return nil
}
