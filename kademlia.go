package d7024e

type Kademlia struct {
	rt   *RoutingTable
	netw *Network
}

func NewKademlia(me Contact, port string) *Kademlia {
	rt := NewRoutingTable(me)
	// hardcoded port for now
	netw := NewNetwork(port, me.Address)
	k := &Kademlia{rt: rt,
		netw: netw}
	return k
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
