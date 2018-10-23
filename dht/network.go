package dht

import (
	pb "D7024E/dht/pb"
	"log"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

type Network struct {
	port             string
	addr             string
	requestMap       map[int32](chan *pb.Message)
	timedoutRequests chan int32
	msgFct           *pb.MessageFactory
	conn             *net.UDPConn
	rMapLck          *sync.Mutex
	stop             chan struct{}
	isListening      bool
}

func NewNetwork(port, addr string) *Network {
	netw := &Network{port: port,
		addr:             addr,
		requestMap:       make(map[int32](chan *pb.Message)),
		msgFct:           pb.NewMessageFactory(),
		timedoutRequests: make(chan int32, 100),
		rMapLck:          &sync.Mutex{},
		stop:             make(chan struct{}),
		isListening:      false}
	return netw
}

func (k *Kademlia) InitConn() {
	n := k.netw
	serverAddr, err := net.ResolveUDPAddr("udp", ":"+n.port)
	if err != nil {
		log.Fatal(err)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	n.conn = serverConn
}

func (k *Kademlia) CloseConn() {
	k.netw.isListening = false
	close(k.netw.stop)
	_ = k.netw.conn.Close()
}

func Listen(k *Kademlia) error {
	network := k.netw
	network.isListening = true
	buf := make([]byte, 4096)

	log.Println("Listening on port " + network.port)
	for {
		select {
		case <-network.stop:
			log.Println("stopping listener")
			return nil
		case reqID := <-k.netw.timedoutRequests:
			network.rMapLck.Lock()
			returnCh, exists := network.requestMap[reqID]
			if exists {
				close(returnCh)
				delete(network.requestMap, reqID)
			}
			network.rMapLck.Unlock()
		default:
			k.read(buf)
		}
	}
}

func (k *Kademlia) read(buf []byte) {
	network := k.netw
	n, addr, err := network.conn.ReadFromUDP(buf)
	if !network.isListening {
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	msg := &pb.Message{}
	err = proto.Unmarshal(buf[0:n], msg)
	if err != nil {
		log.Println(err)
		return
	}
	if msg.Response {
		log.Printf("Received response:\t reqID: %-5d type: %-12s %-6s %s sent at: %s",
			msg.GetRequestID(),
			msg.GetType(),
			"from:",
			addr, time.Unix(msg.GetSentTime(), 0))
		reqID := msg.GetRequestID()
		network.rMapLck.Lock()
		returnCh, exists := network.requestMap[reqID]
		if exists {
			returnCh <- msg
			close(returnCh)
			delete(network.requestMap, reqID)
		}
		network.rMapLck.Unlock()
	} else {
		go k.handleMessage(msg)
		log.Printf("Received request:\t reqID: %-5d type: %-12s %-6s %s sent at: %s",
			msg.GetRequestID(),
			msg.GetType(),
			"from:",
			addr, time.Unix(msg.GetSentTime(), 0))
	}
	go k.updateContacts(msg)
}

func (n *Network) SendRequest(c *Contact, msg *pb.Message,
	returnCh chan *pb.Message) (chan int32, error) {
	err := n.SendMessage(c, msg)
	if err != nil {
		return nil, err
	}
	n.rMapLck.Lock()
	n.requestMap[msg.GetRequestID()] = returnCh
	n.rMapLck.Unlock()
	return n.timedoutRequests, nil
}

func (n *Network) SendMessage(c *Contact,
	msg *pb.Message) error {
	remoteAddr, err := net.ResolveUDPAddr("udp", c.Address)
	if err != nil {
		return err
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = n.conn.WriteToUDP(data, remoteAddr)
	if err != nil {
		return err
	}
	if msg.Response {
		log.Printf("Sent response:\t reqID: %-5d type: %-12s %-6s %s",
			msg.GetRequestID(),
			msg.GetType(),
			"to:",
			c.Address)
	} else {
		log.Printf("Sent request:\t reqID: %-5d type: %-12s %-6s %s",
			msg.GetRequestID(),
			msg.GetType(),
			"to:",
			c.Address)
	}

	return nil
}

func (k *Kademlia) handleMessage(msg *pb.Message) {
	handler := k.getTypeHandler(msg.GetType())
	respMsg, err := handler(msg)
	if err != nil {
		log.Println(err)
		return
	}
	if respMsg == nil {
		return
	}
	receiver := PeerToContact(respMsg.GetReceiver())
	err = k.netw.SendMessage(&receiver, respMsg)
	if err != nil {
		log.Println(err)
	}
}

func (k *Kademlia) updateContacts(msg *pb.Message) {
	//Don't update on ping responses
	//Maybe we want a flag instead that you set for when you dont want
	// to update the routingtable after certain messages
	//For now we can deal with ignoring updates on ping responses to avoid loops
	if msg.GetType() == pb.Message_PING && msg.Response {
		return
	} else {
		k.Update(PeerToContact(msg.GetSender()))
	}

}
