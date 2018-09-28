package dht

import (
	pb "D7024E/dht/pb"
	"log"
	"net"
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
}

func NewNetwork(port, addr string) *Network {
	netw := &Network{port: port,
		addr:             addr,
		requestMap:       make(map[int32](chan *pb.Message)),
		msgFct:           pb.NewMessageFactory(),
		timedoutRequests: make(chan int32)}
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

func Listen(k *Kademlia, ip string, port string) error {
	network := k.netw
	defer network.conn.Close()

	buf := make([]byte, 1024)

	log.Println("Listening on port " + network.port)
	for {
		n, addr, err := network.conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		msg := &pb.Message{}
		err = proto.Unmarshal(buf[0:n], msg)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Received %d sent at %s from %s",
			msg.GetMessageID(),
			time.Unix(msg.GetSentTime(), 0), addr)
		if err != nil {
			log.Fatal(err)
		}
		if msg.Response {
			reqID := msg.GetRequestID()
			returnCh := network.requestMap[reqID]
			returnCh <- msg
			close(returnCh)
			delete(network.requestMap, reqID)
			//TODO: cleanup timed-out messagehandlers
		} else {
			go k.handleMessage(msg)
		}
		go k.updateContacts(msg)
	}
}

func (n *Network) SendRequest(c *Contact, msg *pb.Message,
	returnCh chan *pb.Message) (chan int32, error) {
	n.requestMap[msg.GetRequestID()] = returnCh
	err := n.SendMessage(c, msg)
	if err != nil {
		delete(n.requestMap, msg.GetRequestID())
		return nil, err
	}
	return n.timedoutRequests, nil
}

func (n *Network) SendMessage(c *Contact,
	msg *pb.Message) error {
	remoteAddr, err := net.ResolveUDPAddr("udp", c.Address+":"+n.port)
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
	log.Printf("Sent message %d type %s",
		msg.GetMessageID(), msg.GetType())

	return nil
}

func (k *Kademlia) handleMessage(msg *pb.Message) {
	handler := k.getTypeHandler(msg.GetType())
	respMsg, err := handler(msg)
	if err != nil {
		log.Fatal(err)
	}
	receiver := PeerToContact(respMsg.GetReceiver())
	err = k.netw.SendMessage(&receiver, respMsg)
	if err != nil {
		log.Fatal(err)
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

// Dont use these. Instead, create a message in msgFct and send it using SendMessage
func (n *Network) SendPingMessage(tar, me *Contact, reqID int32) error {
	// TODO
	return nil
}

func (n *Network) SendFindContactMessage(c *Contact) error {
	// TODO
	return nil
}

func (n *Network) SendFindDataMessage(hash string) error {
	// TODO
	return nil
}

func (n *Network) SendStoreMessage(data []byte) error {
	// TODO
	return nil
}
