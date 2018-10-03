package dht

import (
	"fmt"
	"testing"
)

func TestContact(t *testing.T) {
	
	TestcontactToPeer(t)
	TestpeerToContact(t)
	TestcontactsToPeers(t)
	TestpeersToContacts(t)
	
	fmt.Println("anal")
}

func TestcontactToPeer(t *testing.T) {
		
	contact := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "192.169.0.0")
	
	_ = ContactToPeer(contact)
	
	//Sets the distance for that contact
	contact.distance = NewKademliaID("1000000000000000000000000000000000000000")
	
	//Calls the method again but now the distance is set
	_ = ContactToPeer(contact)

}

func TestpeerToContact(t *testing.T) {
	
	contact := NewContact(NewKademliaID("fffffff000000000000000000000000000000000"), "192.169.0.0")
	
	contact.distance = NewKademliaID("1000000000000000000000000000000000000000")
	
	peer := ContactToPeer(contact)
	
	_ = PeerToContact(peer)
}

func TestcontactsToPeers(t *testing.T) {
	
	contacts := make([]Contact, 0)
	
	contact := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "192.169.0.0")
	
	contacts = append(contacts, contact)
	
	contact = NewContact(NewKademliaID("fffffff000000000000000000000000000000000"), "192.169.0.1")
	
	contacts = append(contacts, contact)
	
	_ = ContactsToPeers(contacts)
}

func TestpeersToContacts(t *testing.T) {
	
	contacts := make([]Contact, 0)
	
	contact := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "192.169.0.0")
	
	contacts = append(contacts, contact)
	
	contact = NewContact(NewKademliaID("fffffff000000000000000000000000000000000"), "192.169.0.1")
	
	contacts = append(contacts, contact)
	
	peers := ContactsToPeers(contacts)
	
	_ = PeersToContacts(peers)
	
	
	//This tests the function Add in contact 
	
	candidates := ContactCandidates{contacts}
	
	candidates.Add(contact)
	
	fmt.Println(candidates)
	
	//This tests remove
	
	candidates.Remove(1)
	
	//This populates candidates with 20 identical contacts and then cuts of some so that the len() is 20
	//AKA tests Cut()
	
	for i := 0; i < 20; i++ {
		candidates.Add(contact)
	}
	
	candidates.Cut()
	
	
	//This tests Exists() and AddUnique()
	
	contactUnique := NewContact(NewKademliaID("fff0000000000000000000000000000000000000"), "192.169.0.1")
	
	candidates.Exists(contactUnique)
	
	candidates.AddUnique(contacts)
	
	contacts = append(contacts, contactUnique)
	
	candidates.AddUnique(contacts)

	fmt.Println(contacts)
}











