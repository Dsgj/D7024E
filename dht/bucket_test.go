package dht

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func TestBucket(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	buck1 := newBucket()
	buck2 := newBucket()

	//Testar NeedsRefresh()
	buck1.NeedsRefresh(time.Now())
	buck1.Refresh(time.Now())
	fmt.Println(buck1)
	buck1.AddContact(NewContact(NewKademliaID("fffffff000000000000000000000000000000001"), "192.169.0.0"))
	fmt.Println(buck1)
	buck1.AddContact(NewContact(NewKademliaID("ffffff0000000000000000000000000000000010"), "192.169.0.0"))
	fmt.Println(buck1)
	buck1.AddContact(NewContact(NewKademliaID("fffff00000000000000000000000000000000100"), "192.169.0.0"))
	fmt.Println(buck1)
	//Jag tror inte att denhär kan testas helt, iom att det är random så varierar testprocenten

	fmt.Println(buck1.GetRandomContact())
	fmt.Println(buck2.GetRandomContact())
}