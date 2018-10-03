package dht

import (
	//"fmt"
	"testing"
	"time"
)

func TestBucket(t *testing.T) {
	
	buck1 := newBucket()
	
	//Testar NeedsRefresh()
	buck1.NeedsRefresh(time.Now())
	
	buck1.Refresh(time.Now())
	
	buck1.AddContact(NewContact(NewKademliaID("fffffff000000000000000000000000000000000"), "192.169.0.0"))
	buck1.AddContact(NewContact(NewKademliaID("ffffff0000000000000000000000000000000000"), "192.169.0.0"))
	buck1.AddContact(NewContact(NewKademliaID("fffff00000000000000000000000000000000000"), "192.169.0.0"))
	
	//Jag tror inte att denhär kan testas helt, iom att det är random så varierar testprocenten
	
	_ = buck1.GetRandomContact()

}











