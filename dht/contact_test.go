package dht

import (
	"testing"
)

func TestContactToPeer(t *testing.T) {
	id := "ffffffff00000000000000000000000000000000"
	ip := "192.169.0.0"
	contact := NewContact(NewKademliaID(id), ip)
	peer := ContactToPeer(contact)

	if peer.Id != id {
		t.Errorf("peer.Id was incorrect, got: %s, want: %s.", peer.Id, id)
	}
	if peer.Addr != ip {
		t.Errorf("peer.Addr was incorrect, got: %s, want: %s.", peer.Addr, ip)
	}
	if peer.Distance != "" {
		t.Errorf("peer.Distance was incorrect, got: %s, want: %s.", peer.Distance, "")
	}

	distance := NewRandomKademliaID()
	contact.distance = distance
	peer = ContactToPeer(contact)

	if peer.Distance != distance.String() {
		t.Errorf("peer.Distance was incorrect, got: %s, want: %s.\n", peer.Distance, distance.String())
	}
}

func TestPeerToContact(t *testing.T) {
	id := "ffffffff00000000000000000000000000000000"
	ip := "192.169.0.0"

	contact1 := NewContact(NewKademliaID(id), ip)
	peer := ContactToPeer(contact1)
	contact2 := PeerToContact(peer)
    contact3 := NewContact(NewRandomKademliaID(), "192.169.0.0")
	contact2.CalcDistance(contact1.ID)
	contact1.CalcDistance(contact2.ID)

    contact3.CalcDistance(contact1.ID)

    contact1.String()
	contact2.String()
	contact3.String()


	if !contact1.ID.Equals(contact2.ID) {
		t.Errorf("contact2.ID was incorrect, got: %s, want: %s.\n", contact2.ID, contact1.ID)
	}
	if contact1.Address != contact2.Address {
		t.Errorf("contact2.Address was incorrect, got: %s, want: %s.\n", contact2.Address, contact1.Address)
	}
	if !contact1.distance.Equals(contact2.distance) {
		t.Errorf("contact2.distance was incorrect, got: %s, want: %s.\n", contact2.distance, contact1.distance)
	}
    if !contact2.distance.Less(contact3.distance) {
		t.Errorf("contact2.distance was incorrect, got: %s, want: %s.\n", contact2.distance, contact3.distance)
	}
}

func TestContactsToPeers(t *testing.T) {
	contacts := make([]Contact, 0)
	contact := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "192.169.0.0")
	contacts = append(contacts, contact)
	contact = NewContact(NewKademliaID("fffffff000000000000000000000000000000000"), "192.169.0.1")
	contacts = append(contacts, contact)
	_ = ContactsToPeers(contacts)
}

func TestPeersToContacts(t *testing.T) {
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
	//fmt.Println(candidates)

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
	//fmt.Println(contacts)
}
