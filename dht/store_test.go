package dht

import (
	"fmt"
	"testing"
	"time"
	//"encoding/hex"
)

func TestStore(t *testing.T) {
	Teststore(t)
}

func Teststore(t *testing.T) {
	stor1 := NewStore()
	data := []byte{'X'}
	contact := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "192.169.0.0")

	//Tests store
	_ = stor1.Store(data, contact, time.Now())

	//Test getRecord with both a key that exists and with one that doesnt exist
	key := GetKey(data)
	record, _ := stor1.GetRecord(key)
	data2 := []byte{'X', 'Y', 'Z'}
	keyFake := GetKey(data2)
	stor1.GetRecord(keyFake)

	//Tests PinRecord() with both a real key and a false key(a key to data that is not stored)
	stor1.PinRecord(key)
	stor1.PinRecord(keyFake)

	//Tests DelRecord. Tries to delete a record that does not exist, a record that is pinned and a record that is real and unpinned
	stor1.DelRecord(key)
	stor1.UnpinRecord(key)
	stor1.DelRecord(key)
	stor1.DelRecord(keyFake)

	//Tests UnpinRecord by trying to unpin with a fake key
	stor1.UnpinRecord(keyFake)

	//Tests SendableRecord() by passing both a real and a fake key
	stor1.Store(data, contact, time.Now())
	stor1.SendableRecord(key, false)
	stor1.SendableRecord(keyFake, true)

	//Tests Republish()
	record.Republish(time.Now())

	//Tests Replicate()
	record.Replicate(time.Now())

	//Tests IsExpired()
	record.IsExpired(time.Now())

	//Tests NeedsReplicate()
	record.NeedsReplicate(time.Now())

	//Tests NeedRepublish()
	record.NeedsRepublish(time.Now(), contact)

	Testfromstring(t)
}

func Testfromstring(t *testing.T) {
	//Tests FromString(), the input has to be long enough so that the for-loop doesnt get a outofbounds error with the bytearray[20]
	fmt.Println("len > 20")
	FromString("48656c6c6f20476f706865722148656c6c6f20476f7068657221")
	fmt.Println(FromString("48656c6c6f20476f706865722148656c6c6f20476f7068657221"))
	fmt.Println("\nlen == 20")
	FromString("48656c6c6f20476f706865722148656c6c6f20476")
	fmt.Println(FromString("48656c6c6f20476f706865722148656c6c6f20476"))
	fmt.Println("\nlen < 20")
	FromString("48656c6c6f20476f7068657221")
	fmt.Println(FromString("48656c6c6f20476f7068657221"))
}

