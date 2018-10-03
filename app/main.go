package main

import (
	d "D7024E/dht"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	port          = "8080"
	bootstrapID   = "210fc7bb818639ac48a4c6afa2f1581a8b9525e3"
	bootstrapAddr = "10.0.0.4"
)

func main() {
	// take ip from args for now
	ip := os.Args[1]
	me := getSelf(ip)
	k := d.NewKademlia(me, port)

	// init listener/conn
	k.InitConn()
	go d.Listen(k)

	if ip != bootstrapAddr {
		bootstrap(k, me)
	}
	k.StartScheduler()
	select {} // block forever
}

func bootstrap(k *d.Kademlia, me d.Contact) {
	bs := d.NewContact(d.NewKademliaID(bootstrapID), bootstrapAddr)
	k.Update(bs)
	closestContacts, err := k.IterativeFindNode(me.ID.String())
	if err != nil {
		log.Println(err)
	}
	if len(closestContacts) == 0 {
		log.Printf("Received no contacts\nRestarting bootstrap")
		time.Sleep(time.Second * 5)
		bootstrap(k, me)
	} else {
		fmt.Printf("received contacts: %+v\n", closestContacts)
	}

}

func getSelf(ip string) d.Contact {
	if ip == bootstrapAddr { // im the bootstrapnode (this is bad, but works)
		fmt.Printf("I'm the bootstrap node\n")
		return d.NewContact(d.NewKademliaID(bootstrapID), ip)
	} else {
		id := d.NewRandomKademliaID()
		return d.NewContact(id, ip)
	}
}
