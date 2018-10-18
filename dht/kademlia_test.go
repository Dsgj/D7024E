package dht

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestKademlia(t *testing.T) {
	fmt.Println("TestKademlia: ENTER!")

	//creates a new instance of Kademlia
	kademlia1 := NewKademlia(NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000"), "1337") //NewKademlia( contact, port(string))

	fmt.Println("Set Port: " + kademlia1.netw.port)

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

	fmt.Println("TestKademlia: EXIT!\n ")
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

	index = kademlia1.rt.getBucketIndex(NewKademliaID("ffffffff00000000000000000000000000000000"))

	lenOfBucket1 = kademlia1.rt.buckets[index].Len()

	kademlia1.Update(NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000"))

	lenOfBucket2 = kademlia1.rt.buckets[index].Len()

	if lenOfBucket1 != lenOfBucket2 {
		t.Errorf("Kademlia.Update() did let a node add itself. \n")
	}

}

func TestKademlia_RandomID(t *testing.T) {
	fmt.Println("TestKademlia_RandomID: ENTER!")
	//creates a new instance of Kademlia
	kademlia3 := NewKademlia(NewContact(NewRandomKademliaID(), "localhost:8000"), "1337") //NewKademlia( contact, port(string))
	time.Sleep(time.Millisecond * 2)
	kademlia2 := NewKademlia(NewContact(NewRandomKademliaID(), "localhost:8000"), "1337")
	fmt.Println("Randomized KademliaID: ", kademlia3.rt.me.ID)

	//checks if the port is set
	if kademlia3.netw.port != "1337" {
		t.Errorf("Port was incorrect, got: %s, want: %s.", kademlia3.netw.port, "1337")
	}

	//checks if the ID is set
	if len(*kademlia3.rt.me.ID) != 20 {
		t.Errorf("The lenght of the KademliaID was incorrect, got: %d, want: %d. \n", len(kademlia3.rt.me.ID), 20)
	}

	//checks if the randomized ID is unique
	fmt.Println("Randomized Kademlia3: ", kademlia3.rt.me.ID)
	fmt.Println("Randomized Kademlia2: ", kademlia2.rt.me.ID)
	if *kademlia3.rt.me.ID == *kademlia2.rt.me.ID {
		t.Errorf("The ID has been generated twice")
	}
	fmt.Println("TestKademlia_RandomID: EXIT!\n ")
}

func TestKademlia_STORE(t *testing.T) {
	kademlia4 := NewKademlia(NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000"), "1337") //NewKademlia( contact, port(string))
	rand.Seed(time.Now().UnixNano())
	N := rand.Intn(10)
	testBytes := make([]byte, N)
	for i := 0; i < N; i++ {
		testBytes[i] = 'a' + byte(i%26)
	}
	rec := kademlia4.dataStore.Store(testBytes, kademlia4.rt.me, time.Now())
	fmt.Printf("iterativestore on rec: %v", rec)
	kademlia4.IterativeStore(GetKey(testBytes), true)
}

func TestKademlia_FIND_VALUE(t *testing.T) {

}
func TestKademlia_FIND_NODE(t *testing.T) {

}
func TestKademlia_PING(t *testing.T) {

}
func TestKademlia_FetchFile(t *testing.T) {

}
func TestKademlia_IterativeFindNode(t *testing.T) {

}
func TestKademlia_IterativeFindValue(t *testing.T) {

}
func TestKademlia_IterativeStore(t *testing.T) {

}
func TestKademlia_StartScheduler(t *testing.T) {

}

func TestPingExample(t *testing.T) {
	contacts, kademlias := InitKademlias(20)
	respMsg, err, timeout := kademlias[0].PING(contacts[5])
	if err != nil {
		fmt.Println(err)
	}
	if timeout {
		fmt.Println("ping timed out")
	}
	fmt.Printf("ping response msg: %v\n", respMsg)

	//do something with message, compare sender etc
}

func InitKademlias(num int) ([]Contact, []*Kademlia) {
	if num < 1 {
		return nil, nil
	}
	contacts := make([]Contact, 0)
	for i := 0; i < num; i++ {
		id := NewRandomKademliaID()
		contacts = append(contacts, NewContact(id, "localhost:100"+fmt.Sprintf("%d", i)))
		time.Sleep(time.Millisecond * 2)
	}
	kademlias := make([]*Kademlia, 0)
	for i := 0; i < num; i++ {
		k := NewKademlia(contacts[0], "100"+fmt.Sprintf("%d", i))
		k.InitConn()
		go Listen(k)
		for _, c := range contacts {
			k.Update(c)
		}
		kademlias = append(kademlias, k)
	}
	return contacts, kademlias
}
