package main

import (
	d "D7024E"
)

const (
	ip   = "127.0.0.1"
	port = "8080"
)

func main() {
	id := d.NewRandomKademliaID()
	me := d.NewContact(id, ip)
	k := d.NewKademlia(me, port)

	d.Listen(k, ip, port)
}
