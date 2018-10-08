package dht

import (
	pb "D7024E/dht/pb"
	"fmt"
	"sort"
)

// Contact definition
// stores the KademliaID, the ip address and the distance
type Contact struct {
	ID       *KademliaID
	Address  string
	distance *KademliaID
}

// NewContact returns a new instance of a Contact
func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil}
}

func ContactToPeer(c Contact) *pb.Peer {
	dist := ""
	if c.distance != nil {
		dist = c.distance.String()
	}
	return &pb.Peer{
		Id:       c.ID.String(),
		Addr:     c.Address,
		Distance: dist,
	}
}

func PeerToContact(p *pb.Peer) Contact {
	var dist *KademliaID
	if p.GetDistance() != "" {
		dist = NewKademliaID(p.GetDistance())
	}
	return Contact{
		ID:       NewKademliaID(p.GetId()),
		Address:  p.GetAddr(),
		distance: dist,
	}
}

func ContactsToPeers(contacts []Contact) []*pb.Peer {
	peers := make([]*pb.Peer, 0)
	for _, contact := range contacts {
		peers = append(peers, ContactToPeer(contact))
	}
	return peers
}

func PeersToContacts(peers []*pb.Peer) []Contact {
	contacts := make([]Contact, 0)
	for _, peer := range peers {
		contacts = append(contacts, PeerToContact(peer))
	}
	return contacts
}

// CalcDistance calculates the distance to the target and
// fills the contacts distance field
func (contact *Contact) CalcDistance(target *KademliaID) {
	contact.distance = contact.ID.CalcDistance(target)
}

// Less returns true if contact.distance < otherContact.distance
func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.distance.Less(otherContact.distance)
}

// String returns a simple string representation of a Contact
func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

func (contact *Contact) Equals(otherContact *Contact) bool {
	return contact.ID.String() == otherContact.ID.String()
}

// ContactCandidates definition
// stores an array of Contacts
type ContactCandidates struct {
	contacts []Contact
}

// Append an array of Contacts to the ContactCandidates
func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}

func (candidates *ContactCandidates) Add(contact Contact) {
	candidates.contacts = append(candidates.contacts, contact)
}

// GetContacts returns the first count number of Contacts
func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	return candidates.contacts[:count]
}

// Sort the Contacts in ContactCandidates
func (candidates *ContactCandidates) Sort() {
	if candidates.Len() > 1 {
		sort.Sort(candidates)
	}
}

// Len returns the length of the ContactCandidates
func (candidates *ContactCandidates) Len() int {
	return len(candidates.contacts)
}

// Swap the position of the Contacts at i and j
// WARNING does not check if either i or j is within range
func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.contacts[i], candidates.contacts[j] = candidates.contacts[j], candidates.contacts[i]
}

// Less returns true if the Contact at index i is smaller than
// the Contact at index j
func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.contacts[i].Less(&candidates.contacts[j])
}

func (candidates *ContactCandidates) Remove(i int) {
	candidates.contacts = append(candidates.contacts[:i], candidates.contacts[i+1:]...)
}

func (candidates *ContactCandidates) Cut() {
	if candidates.Len() > 20 {
		candidates.contacts = candidates.contacts[:20]
	}
}

func (candidates *ContactCandidates) Exists(c Contact) bool {
	for _, contact := range candidates.contacts {
		if contact.Equals(&c) {
			return true
		}
	}
	return false
}

func (candidates *ContactCandidates) AddUnique(contacts []Contact) {
	for _, contact := range contacts {
		if !candidates.Exists(contact) {
			candidates.Add(contact)
		}
	}
}
