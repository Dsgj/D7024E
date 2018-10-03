package dht

import (
	//"fmt"
	"testing"
	pb "D7024E/dht/pb"
)

func TestRpcHandler(t *testing.T) {
	

	kademlia1 := NewKademlia(NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000"), "1337") //NewKademlia( contact, port(string))
	
	//Tests getTypeHandler() by getting all 4 different handlers
	
	kademlia1.getTypeHandler(0)
	kademlia1.getTypeHandler(1)
	kademlia1.getTypeHandler(2)
	kademlia1.getTypeHandler(3)
	
	
	_ = pb.NewMessageFactory()
	
}
































