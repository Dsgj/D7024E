package main

import (
	d "D7024E"
	"log"
	"net"
)

const (
	ip   = "127.0.0.1"
	port = "8080"
)

func main() {
	id := d.NewRandomKademliaID()
	me := d.NewContact(id, ip)
	k := d.NewKademlia(me, port)

	//myip := GetOutboundIP()
	//log.Printf("IP Address: %d", myip)
	//ListIPs()

	go d.Listen(k, ip, port)

	log.Printf("ID: %d\nAddr: %s", id.String(), ip)
	select {} // block forever
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
