package dht

import (
	pb "D7024E/dht/pb"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"sync"
	"time"
)

const ExpirationPeriod = time.Hour * 24
const ReplicationPeriod = time.Hour * 1

type Store struct {
	records map[[20]byte]*Record
	mutex   *sync.Mutex
}

func NewStore() *Store {
	return &Store{records: make(map[[20]byte]*Record),
		mutex: &sync.Mutex{}}
}

type Record struct {
	value    []byte
	expTime  time.Time
	replTime time.Time
	pinned   bool
}

func (s *Store) Store(data []byte) *Record {
	sha := Hash(data)
	record := &Record{value: data,
		expTime:  time.Now(),
		replTime: time.Now(),
		pinned:   false}
	s.mutex.Lock()
	s.records[sha] = record
	s.mutex.Unlock()
	return record
}

func (s *Store) GetRecord(key [20]byte) (*Record, bool) {
	s.mutex.Lock()
	record, exists := s.records[key]
	s.mutex.Unlock()
	if exists {
		return record, true
	} else {
		log.Printf("No record found for key: %s", ToString(key))
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

func (s *Store) PinRecord(key [20]byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if record, exists := s.records[key]; exists {
		record.pinned = true
		log.Printf("Pinned record: %s", ToString(key))
	} else {
		log.Printf("No record found for key: %s", ToString(key))
	}
}

func (s *Store) UnpinRecord(key [20]byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if record, exists := s.records[key]; exists {
		record.pinned = false
		log.Printf("Unpinned record: %s", ToString(key))
	} else {
		log.Printf("No record found for key: %s", ToString(key))
	}
}

func (s *Store) SendableRecord(key [20]byte) *pb.Record {
	rec, exists := s.GetRecord(key)
	if exists {
		return &pb.Record{Key: key[:], Value: rec.value}
	}
	return nil
}

func (r *Record) IsExpired(now time.Time) bool {
	return now.Sub(r.expTime) >= ExpirationPeriod
}

func (r *Record) NeedsReplicate(now time.Time) bool {
	return now.Sub(r.replTime) >= ReplicationPeriod
}

func Hash(data []byte) [20]byte {
	return sha1.Sum(data)
}

func ToString(data [20]byte) string {
	return hex.EncodeToString(data[:])
}

func FromString(s string) [20]byte {
	decoded, _ := hex.DecodeString(s)

	newArray := [20]byte{}
	for i := 0; i < 20; i++ {
		newArray[i] = decoded[i]
	}

	return newArray
}
