package dht

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestKademlia(t *testing.T) {

	//creates a new instance of Kademlia
	kademlia1 := NewKademlia(NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000"), "1337") //NewKademlia( contact, port(string))

	//fmt.Println(kademlia1.netw.port)

	//checks if the port is set
	if kademlia1.netw.port != "1337" {
		t.Errorf("Port was incorrect, got: %s, want: %s.", kademlia1.netw.port, "1337")
	}

	fmt.Println(" ")

	//Creates a test Kademlia ID to compare with the ID from the created instance of Kademlia
	testID := KademliaID{}
	decoded, _ := hex.DecodeString("ffffffff00000000000000000000000000000000")
	for i := 0; i < 20; i++ {
		testID[i] = decoded[i]
	}

	//checks if the ID is set
	if *kademlia1.rt.me.ID != testID {
		t.Errorf("KademliaID was incorrect, got: %d, want: %d.", kademlia1.rt.me.ID, &testID)
	}

	//This checks that the values has ben set
	//---------------------------------------------------

	TestnewRequest(t, *kademlia1)

	Testupdate(t, *kademlia1)

	//test, _ := kademlia1.IterativeFindNode("ffffffff00000000000000000000000000000000")

	//fmt.Println(test)

	//kademlia1.Update( NewContact(NewKademliaID("fffffff000000000000000000000000000000000"), "localhost:8000"))

}

func TestnewRequest(t *testing.T, kademlia1 Kademlia) {
	requestID := kademlia1.newRequest()

	if kademlia1.reqCount != requestID+1 {
		t.Errorf("Request count was incorrect, is: %d, should be: %d.", kademlia1.reqCount, requestID+1)
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
		t.Errorf("Kademlia.Update() did not change the number of contacts in the updated bucket.")
	}

	/* Checks if kademlia.Update lets a node att itself by checking the difference in the size of the bucked that it would have been in
	 */

	index = kademlia1.rt.getBucketIndex(NewKademliaID("ffffffff00000000000000000000000000000000"))

	lenOfBucket1 = kademlia1.rt.buckets[index].Len()

	kademlia1.Update(NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000"))

	lenOfBucket2 = kademlia1.rt.buckets[index].Len()

	if lenOfBucket1 != lenOfBucket2 {
		t.Errorf("Kademlia.Update() did let a node add itself.")
	}

}
