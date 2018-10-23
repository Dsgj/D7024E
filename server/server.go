package main

import (
	d "D7024E/dht"
	"D7024E/server/api"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/takama/daemon"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	port          = "8080"
	bootstrapID   = "210fc7bb818639ac48a4c6afa2f1581a8b9525e3"
	bootstrapAddr = "10.0.0.4"
	bootstrapPort = "8060"
)

func main() {
	service, err := daemon.New("dht server", "dht server daemon with rest api")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	status, err := service.Install()
	if err != nil {
		log.Fatal(status, "\nError: ", err)
	}
	fmt.Println(status)
	// take ip from args for now
	ip := os.Args[1]
	me := getSelf(ip, port)
	var k *d.Kademlia
	if ip != bootstrapAddr {
		k = d.NewKademlia(me, port)
	} else {
		k = d.NewKademlia(me, bootstrapPort)
	}

	// init listener/conn
	k.InitConn()
	go d.Listen(k)

	if ip != bootstrapAddr {
		bootstrap(k, me)
		k.StartScheduler()
	}
	r := chi.NewRouter()
	api.Routes(r, k)
	http.ListenAndServe(":3000", r)
}

func bootstrap(k *d.Kademlia, me d.Contact) {
	bs := d.NewContact(d.NewKademliaID(bootstrapID), bootstrapAddr+":"+bootstrapPort)
	k.Update(bs)
	closestContacts, err := k.IterativeFindNode(me.ID.String())
	if err != nil {
		log.Println(err)
	}
	if len(closestContacts) == 0 {
		log.Printf("Received no contacts, restarting bootstrap")
		time.Sleep(time.Second * 5)
		bootstrap(k, me)
	} else {
		log.Printf("Bootstrap done, received contacts: %+v\n", closestContacts)
		// test store
		k.TestStore()
	}

}

func getSelf(ip, port string) d.Contact {
	if ip == bootstrapAddr { // im the bootstrapnode (this is bad, but works)
		log.Printf("I'm the bootstrap node\n")
		return d.NewContact(d.NewKademliaID(bootstrapID), ip+":"+bootstrapPort)
	} else {
		id := d.NewRandomKademliaID()
		return d.NewContact(id, ip+":"+port)
	}
}
