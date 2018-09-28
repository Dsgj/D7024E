package dht

import (
	"crypto/sha1"
	"encoding/hex"
	"time"
)

const ExpirationPeriod = time.Hour * 24
const ReplicationPeriod = time.Hour * 1

type Store struct {
	records map[[20]byte]*Record
}

type Record struct {
	value    []byte
	expTime  time.Time
	replTime time.Time
}

func (s *Store) Store(data []byte) {
	sha := Hash(data)
	record := &Record{value: data,
		expTime:  time.Now(),
		replTime: time.Now()}
	s.records[sha] = record
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

func FromString(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
