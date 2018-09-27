package main

import (
	d "D7024E/dht"
	"fmt"
	"log"
	"net"
	"os"
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

	//myip := GetOutboundIP()
	//log.Printf("IP Address: %d", myip)
	//ListIPs()

	// init listener/conn
	k.InitConn()
	//TODO: listen does not use these params for now, clean up
	go d.Listen(k, ip, port)

	if ip != bootstrapAddr {
		bootstrap(k, me)
	}
	select {} // block forever
}

func bootstrap(k *d.Kademlia, me d.Contact) {
	bs := d.NewContact(d.NewKademliaID(bootstrapID), bootstrapAddr)
	closestContacts, err, timeout := k.FIND_NODE(bs, me.ID.String())
	if err != nil {
		log.Fatal(err)
	}
	if timeout {
		log.Printf("Request timed out\nRestarting bootstrap")
		bootstrap(k, me)
	}
	fmt.Printf("%+v\n", closestContacts)
	for _, contact := range closestContacts {
		k.Update(contact)
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

// gets preferred outbound IP
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// prints all IPs
func ListIPs() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	count := 1
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			log.Printf("IP Address %d: %d", count, ip)
			count++
		}
	}
}
