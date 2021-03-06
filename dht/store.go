package dht

import (
	pb "D7024E/dht/pb"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// const tExpire = time.Hour * 25
// const tRepublish = time.Hour * 24
// const tReplicate = time.Hour * 1
const tExpire = time.Second * 60    //testing
const tRepublish = time.Second * 45 //testing
const tReplicate = time.Second * 20 //testing

type Store struct {
	records map[[20]byte]*Record
	mutex   *sync.Mutex
}

func NewStore() *Store {
	return &Store{records: make(map[[20]byte]*Record),
		mutex: &sync.Mutex{}}
}

type Record struct {
	value       []byte
	replAt      time.Time
	publishedAt time.Time
	publisher   Contact
	pinned      bool
	mutex       *sync.Mutex
}

func (s *Store) Store(data []byte, publisher Contact, publAt time.Time) *Record {
	sha := GetKey(data)
	record := &Record{value: data,
		replAt:      time.Now(),
		publishedAt: publAt,
		publisher:   publisher,
		pinned:      false,
		mutex:       s.mutex}
	s.mutex.Lock()
	s.records[sha] = record
	s.mutex.Unlock()
	log.Printf("record added to store: %v", record)
	return record
}

func (s *Store) GetRecord(key [20]byte) (*Record, bool) {
	s.mutex.Lock()
	record, exists := s.records[key]
	s.mutex.Unlock()
	if exists {
		return record, true
	} else {
		return nil, false
	}

}

func (s *Store) DelRecord(key [20]byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if record, exists := s.records[key]; exists {
		if record.pinned {
			log.Printf("Attempted deletion of pinned record: %s", ToString(key))
		} else {
			delete(s.records, key)
			log.Printf("Deleted record: %s", ToString(key))
		}
	} else {
		log.Printf("No record found for key: %s", ToString(key))
	}
}

func (s *Store) PinRecord(key [20]byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if record, exists := s.records[key]; exists {
		record.pinned = true
		return nil
	}
	return fmt.Errorf("No record found for key: %s", ToString(key))
}

func (s *Store) UnpinRecord(key [20]byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if record, exists := s.records[key]; exists {
		record.pinned = false
		return nil
	}
	return fmt.Errorf("No record found for key: %s", ToString(key))
}

func (s *Store) SendableRecord(key [20]byte, newPublish bool) *pb.Record {
	rec, exists := s.GetRecord(key)
	if exists {
		return &pb.Record{Key: key[:],
			Value:       rec.value,
			NewPublish:  newPublish,
			Publisher:   ContactToPeer(rec.publisher),
			PublishedAt: rec.publishedAt.Unix()}
	}
	return nil
}

func (r *Record) Republish(t time.Time) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.publishedAt = t
	r.replAt = t
	log.Printf("record publishedAt updated: %v", r)
}

func (r *Record) Replicate(t time.Time) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.replAt = t
	log.Printf("record replAt updated: %v", r)
}

func (r *Record) IsExpired(now time.Time) bool {
	return now.Sub(r.publishedAt) >= tExpire && r.pinned == false
}

func (r *Record) NeedsReplicate(now time.Time) bool {
	return now.Sub(r.replAt) >= tReplicate
}

func (r *Record) NeedsRepublish(now time.Time, me Contact) bool {
	return now.Sub(r.publishedAt) >= tRepublish && me.Equals(&r.publisher)
}

func GetKey(data []byte) [20]byte {
	return sha1.Sum(data)
}

func ToString(data [20]byte) string {
	return hex.EncodeToString(data[:])
}

func FromString(s string) ([20]byte, error) {
	decoded, _ := hex.DecodeString(s)
	newArray := [20]byte{}
	if len(decoded) != 20 {
		return newArray, errors.New("cant convert string, length must equal 20")
	}
	for i := 0; i < 20; i++ {
		newArray[i] = decoded[i]
	}

	return newArray, nil
}
