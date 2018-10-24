package dht

import (
	"D7024E/dht/pb"
	"encoding/hex"
	"fmt"
	"io/ioutil" //toggle log show/hide
	"log"       //toggle log show/hide
	"math/rand"
	"testing"
	"time"
	//"reflect"
)

func TestKademlia(t *testing.T) {
	log.SetOutput(ioutil.Discard) //toggle log show/hide

	//creates a new instance of Kademlia
	kademlia1 := NewKademlia(NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000"), "1337") //NewKademlia( contact, port(string))

	//checks if the port is set
	if kademlia1.netw.port != "1337" {
		t.Errorf("Port was incorrect, got: %s, want: %s.", kademlia1.netw.port, "1337")
	}

	//Creates a test Kademlia ID to compare with the ID from the created instance of Kademlia
	testID := KademliaID{}
	decoded, _ := hex.DecodeString("ffffffff00000000000000000000000000000000")
	for i := 0; i < 20; i++ {
		testID[i] = decoded[i]
	}

	//checks if the ID is set
	if *kademlia1.rt.me.ID != testID {
		t.Errorf("KademliaID was incorrect, got: %d, want: %d. \n", kademlia1.rt.me.ID, &testID)
	}

	//This checks that the values has ben set
	//---------------------------------------------------

	TestnewRequest(t, *kademlia1)

	Testupdate(t, *kademlia1)
}
func TestnewRequest(t *testing.T, kademlia1 Kademlia) {
	requestID := kademlia1.newRequest()
	if kademlia1.reqCount != requestID+1 {
		t.Errorf("Request count was incorrect, is: %d, should be: %d.\n", kademlia1.reqCount, requestID+1)
	}
}
func Testupdate(t *testing.T, kademlia1 Kademlia) {

	/* Checks at witch bucket index a kademliaID would be stored at, and saves the length of that bucket.
	 * It then updates with that ID and checks if the lengt of that bucket has changed.
	 */

	index := kademlia1.rt.getBucketIndex(NewKademliaID("fff0000000000000000000000000000000000000"))
	lenOfBucket1 := kademlia1.rt.buckets[index].Len()

	kademlia1.Update(NewContact(NewKademliaID("fff0000000000000000000000000000000000000"), "localhost:8000"))
	lenOfBucket2 := kademlia1.rt.buckets[index].Len()

	if lenOfBucket1 == lenOfBucket2 {
		t.Errorf("Kademlia.Update() did not change the number of contacts in the updated bucket. \n")
	}

	/* Checks if kademlia.Update lets a node att itself by checking the difference in the size of the bucked that it would have been in
	 */

	index2 := kademlia1.rt.getBucketIndex(NewKademliaID("ffffffff00000000000000000000000000000000"))
	//fmt.Println(index2)
	lenOfBucket3 := kademlia1.rt.buckets[index2].Len()
	kademlia1.Update(NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000"))
	lenOfBucket4 := kademlia1.rt.buckets[index2].Len()

	if lenOfBucket3 != lenOfBucket4 {
		t.Errorf("Kademlia.Update() did let a node add itself. \n")
	}
}
func TestKademlia_RandomID(t *testing.T) {
	//creates a new instance of Kademlia
	kademlia3 := NewKademlia(NewContact(NewRandomKademliaID(), "localhost:8000"), "1337") //NewKademlia( contact, port(string))
	time.Sleep(time.Millisecond * 2)
	kademlia2 := NewKademlia(NewContact(NewRandomKademliaID(), "localhost:8000"), "1337")

	//checks if the port is set
	if kademlia3.netw.port != "1337" {
		t.Errorf("Port was incorrect, got: %s, want: %s.", kademlia3.netw.port, "1337")
	}

	//checks if the ID is set
	if len(*kademlia3.rt.me.ID) != 20 {
		t.Errorf("The lenght of the KademliaID was incorrect, got: %d, want: %d. \n", len(kademlia3.rt.me.ID), 20)
	}

	//checks if the randomized ID is unique
	if *kademlia3.rt.me.ID == *kademlia2.rt.me.ID {
		t.Errorf("The ID has been generated twice")
	}
}
func TestKademlia_STORE(t *testing.T) {

	contacts, kademlias, teardown := setupTestCase(t, 20)

	defer teardown(t)

	rand.Seed(time.Now().UnixNano())
	N := rand.Intn(9)
	N += 1
	testBytes := make([]byte, N)
	for i := 0; i < N; i++ {
		testBytes[i] = 'a' + byte(i%26)
	}

	rec := kademlias[0].dataStore.Store(testBytes, kademlias[0].rt.me, time.Now())

	err := kademlias[0].STORE(contacts[1], rec, true)

	if err != nil {
		t.Errorf("STORE genereted an error: %d.\n", err)
	}

	time.Sleep(time.Millisecond * 2)

	value, _, err, timeout := kademlias[0].FIND_VALUE(contacts[1], ToString(GetKey(testBytes)))

	if timeout {
		t.Errorf("FIND_VALUE in test function for STORE got a timeout!")
	}
	if err != nil {
		t.Errorf("FIND_VALUE in test function for STORE generated error: %d.\n", err)
	}

	for i := 0; i < len(value); i++ {
		if value[i] != testBytes[i] {
			t.Errorf("STORE stored the wrong value! Wanted to store: %d. Stored: %d.\n", testBytes, value)
		}
	}

}

func TestKademlia_FIND_VALUE(t *testing.T) {

	contacts, kademlias, teardown := setupTestCase(t, 20)

	defer teardown(t)

	rand.Seed(time.Now().UnixNano())
	N := rand.Intn(9)
	N += 1
	testValue := make([]byte, N)
	for i := 0; i < N; i++ {
		testValue[i] = 'a' + byte(i%26)
	}

	rec := kademlias[0].dataStore.Store(testValue, kademlias[0].rt.me, time.Now())

	err := kademlias[0].STORE(contacts[1], rec, true)

	if err != nil {
		t.Errorf("STORE in test for FIND_NODE genereted an error: %d.\n", err)
	}

	time.Sleep(time.Millisecond * 2)

	//Tries to get the value from a different node than the one that stored it
	_, closestContacts, _, _ := kademlias[0].FIND_VALUE(contacts[2], ToString(GetKey(testValue)))

	if closestContacts == nil {
		t.Errorf("FIND_VALUE did not return closest contacts.")
	}

	//Tries to get the value from the node that we stored it in
	value, _, err, timeout := kademlias[0].FIND_VALUE(contacts[1], ToString(GetKey(testValue)))

	if timeout {
		t.Errorf("FIND_VALUE in test function for FIND_VALUE got a timeout!")
	}
	if err != nil {
		t.Errorf("FIND_VALUE in test function for FIND_VALUE generated error: %d.\n", err)
	}

	for i := 0; i < len(value); i++ {
		if value[i] != testValue[i] {
			t.Errorf("FIND_VALUE returned the wrong value! Wanted: %d. Got: %d.\n", testValue, value)
		}
	}

}
func TestKademlia_FIND_NODE(t *testing.T) {

	contacts, kademlias, teardown := setupTestCase(t, 20)

	defer teardown(t)

	respMsg, err, timeout := kademlias[0].FIND_NODE(contacts[10], contacts[19].ID.String())

	if err != nil {
		t.Errorf("FIND_NODE had an error: %d.\n", err)
	}
	if timeout {
		t.Errorf("FIND_NODE got timeout!")
	}

	externContacts := respMsg
	internContacts := kademlias[10].rt.FindClosestContacts(contacts[19].ID, 20, contacts[0])

	for i := 0; i < 19; i++ {
		if !externContacts[i].ID.Equals(internContacts[i].ID) {
			t.Errorf("FIND_NODE was incorrect, got: %d, want: %d. \n", externContacts[i].ID, internContacts[i].ID)
		}
	}
}
func TestKademlia_PING(t *testing.T) {

}
func TestKademlia_FetchFile(t *testing.T) {

	_, kademlias, teardown := setupTestCase(t, 2)

	defer teardown(t)

	rand.Seed(time.Now().UnixNano())
	N := rand.Intn(9)
	N += 1
	testValue := make([]byte, N)
	for i := 0; i < N; i++ {
		testValue[i] = 'a' + byte(i%26)
	}

	kademlias[0].dataStore.Store(testValue, kademlias[0].rt.me, time.Now())

	file := kademlias[0].FetchFile(ToString(GetKey(testValue)))

	if file == nil {
		t.Errorf("FetchFile did not find the file!")
	}

	file2 := kademlias[1].FetchFile(ToString(GetKey(testValue)))

	if file2 != nil {
		t.Errorf("FetchFile did generate an error or found a file that was stored: %d.\n", file2)
	}

}
func TestKademlia_IterativeFindNode(t *testing.T) {
	numNodes := 20
	_, kademlias, teardown := setupTestCase(t, numNodes)

	defer teardown(t)

	closestContacts, err := kademlias[0].IterativeFindNode(kademlias[2].rt.me.ID.String())
	if err != nil {
		t.Error(err)
	}
	if len(closestContacts) != numNodes-1 {
		t.Errorf("number of contacts received incorrect, expected: %v, got: %v \n",
			numNodes-1, len(closestContacts))
	}
	sorted := &ContactCandidates{}
	sorted.Append(closestContacts)
	sorted.Sort()
	for i, _ := range closestContacts {
		if sorted.contacts[i] != closestContacts[i] {
			t.Errorf("list of contacts was not sorted")
		}
	}

}
func TestKademlia_IterativeFindValue(t *testing.T) {
	numNodes := 20
	_, kademlias, teardown := setupTestCase(t, numNodes)

	defer teardown(t)

	value, _, err := kademlias[0].IterativeFindValue(kademlias[2].rt.me.ID.String())
	if err != nil {
		t.Error(err)
	}
	if value != nil {
		t.Errorf("got value, expected nil")
	}
	rand.Seed(time.Now().UnixNano())
	N := rand.Intn(9)
	N += 1
	testValue := make([]byte, N)
	for i := 0; i < N; i++ {
		testValue[i] = 'a' + byte(i%26)
	}

	rec := kademlias[10].dataStore.Store(testValue, kademlias[10].rt.me, time.Now())
	key := ToString(GetKey(rec.value))
	time.Sleep(time.Millisecond * 5)
	value, _, err = kademlias[0].IterativeFindValue(key)
	if err != nil {
		t.Error(err)
	}
	actual := key
	expected := ToString(GetKey(rec.value))
	if actual != expected {
		t.Errorf("value incorrect, expected: %v, got: %v \n",
			actual, expected)
	}

}
func TestKademlia_IterativeStore(t *testing.T) {
	numNodes := 20
	_, kademlias, teardown := setupTestCase(t, numNodes)

	defer teardown(t)

	rand.Seed(time.Now().UnixNano())
	N := rand.Intn(9)
	N += 1
	testValue := make([]byte, N)
	for i := 0; i < N; i++ {
		testValue[i] = 'a' + byte(i%26)
	}

	rec := kademlias[0].dataStore.Store(testValue, kademlias[0].rt.me, time.Now())
	key := GetKey(rec.value)
	kademlias[0].IterativeStore(key, true) //publish = true

	rec, exists := kademlias[0].dataStore.GetRecord(key)
	if !exists {
		t.Errorf("Record not found, should exist")
	}
	actual := ToString(key)
	expected := ToString(GetKey(rec.value))
	if actual != expected {
		t.Errorf("value incorrect, expected: %v, got: %v \n",
			actual, expected)
	}

}
func TestKademlia_StartScheduler(t *testing.T) {

}
func TestKademlia_PinFile(t *testing.T) {

	_, kademlias, teardown := setupTestCase(t, 2)

	defer teardown(t)

	rand.Seed(time.Now().UnixNano())
	N := rand.Intn(9)
	N += 1
	testValue := make([]byte, N)
	for i := 0; i < N; i++ {
		testValue[i] = 'a' + byte(i%26)
	}

	rec := kademlias[0].dataStore.Store(testValue, kademlias[0].rt.me, time.Now())

	kademlias[0].PinFile(ToString(GetKey(testValue)))

	if !rec.pinned {
		t.Errorf("PinFile did not pin the file!")
	}

	err := kademlias[1].PinFile(ToString(GetKey(testValue)))

	if err == nil {
		t.Errorf("PinFile did not generate an error when trying to pin a file that was not stored!")
	}

}

func TestKademlia_UnpinFile(t *testing.T) {

	_, kademlias, teardown := setupTestCase(t, 2)

	defer teardown(t)

	rand.Seed(time.Now().UnixNano())
	N := rand.Intn(9)
	N += 1
	testValue := make([]byte, N)
	for i := 0; i < N; i++ {
		testValue[i] = 'a' + byte(i%26)
	}

	rec := kademlias[0].dataStore.Store(testValue, kademlias[0].rt.me, time.Now())

	kademlias[0].PinFile(ToString(GetKey(testValue)))

	if rec.pinned {

		kademlias[0].UnpinFile(ToString(GetKey(testValue)))

		if rec.pinned {
			t.Errorf("UnpinFile file did not unpin the file!")
		}
	} else {
		t.Errorf("PinFile in the test for Unpinfile did not work!")
	}

	err := kademlias[1].UnpinFile(ToString(GetKey(testValue)))

	if err == nil {
		t.Errorf("UnpinFile did not generate an error when trying to unpin a file that was not stored!")
	}

}

func TestPingRPC(t *testing.T) {
	contacts, kademlias, teardown := setupTestCase(t, 20)
	defer teardown(t)
	respMsg, err, timeout := kademlias[0].PING(contacts[5])
	if err != nil {
		t.Error(err)
	}
	if timeout {
		t.Errorf("ping timed out")
	}
	//fmt.Printf("ping response msg: %v\n", respMsg)

	sender := PeerToContact(respMsg.GetSender())
	if !sender.Equals(&contacts[5]) {
		t.Errorf("sender of ping response is incorrect, expected: %v, got: %v \n",
			contacts[5], sender)
	}
	if respMsg.GetType() != pb.Message_PING {
		t.Errorf("type of message is incorrect, expected: %s, got: %s \n",
			pb.Message_PING, respMsg.GetType())
	}
}
func InitKademlias(num int) ([]Contact, []*Kademlia) {
	if num < 1 {
		return nil, nil
	}
	contacts := make([]Contact, 0)
	for i := 0; i < num; i++ {
		id := NewRandomKademliaID()
		contacts = append(contacts, NewContact(id, "localhost:500"+fmt.Sprintf("%d", i)))
		time.Sleep(time.Millisecond * 2)
	}
	kademlias := make([]*Kademlia, 0)
	for i := 0; i < num; i++ {
		k := NewKademlia(contacts[0], "500"+fmt.Sprintf("%d", i))
		k.InitConn()
		go Listen(k)
		for _, c := range contacts {
			k.Update(c)
		}
		kademlias = append(kademlias, k)
	}
	return contacts, kademlias
}
func setupTestCase(t *testing.T, num int) ([]Contact, []*Kademlia, func(t *testing.T)) {
	//t.Log("setup test case")
	c, k := InitKademlias(num)
	return c, k, func(t *testing.T) {
		//t.Log("teardown test case")
		for _, kad := range k {
			kad.CloseConn()
		}
	}
}
