package d7024e

import (
	pb "D7024E/pb"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

type Network struct {
	port        string
	addr        string
	msgHandlers map[int32]*MessageHandler
	msgFct      *pb.MessageFactory
	conn        *net.UDPConn
}

func NewNetwork(port, addr string) *Network {
	netw := &Network{port: port,
		addr:        addr,
		msgHandlers: make(map[int32]*MessageHandler),
		msgFct:      pb.NewMessageFactory()}
	return netw
}

type MessageHandler struct {
	id int32
	ch chan (*pb.Message)
}

func (mh *MessageHandler) awaitMessage(returnChan chan (*pb.Message), timeoutChan chan int32) {
	defer close(mh.ch)
	select {
	case result := <-mh.ch:
		returnChan <- result
		return
	case <-time.After(30 * time.Second):
		//TODO: cleanup timed-out messagehandlers
		timeoutChan <- mh.id
		return
	}
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
			returnCh := network.msgHandlers[reqID].ch
			returnCh <- msg
			delete(network.msgHandlers, reqID)
			//TODO: cleanup timed-out messagehandlers
		} else {
			go k.handleMessage(msg)
		}
		go k.updateContacts(msg)
	}
}

func (n *Network) SendMessage(c *Contact,
	msg *pb.Message,
	wantResponse bool) (*MessageHandler, error) {
	remoteAddr, err := net.ResolveUDPAddr("udp", c.Address+":"+n.port)
	if err != nil {
		return nil, err
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	_, err = n.conn.WriteToUDP(data, remoteAddr)
	if err != nil {
		return nil, err
	}
	log.Printf("Sent message %d type %s",
		msg.GetMessageID(), msg.GetType())

	if wantResponse {
		msgHandler := &MessageHandler{
			id: msg.GetRequestID(),
			ch: make(chan (*pb.Message)),
		}
		n.msgHandlers[msgHandler.id] = msgHandler
		return msgHandler, nil
	}

	return nil, nil
}

func (k *Kademlia) handleMessage(msg *pb.Message) {
	handler := k.getTypeHandler(msg.GetType())
	respMsg, err := handler(msg)
	if err != nil {
		log.Fatal(err)
	}
	receiver := PeerToContact(respMsg.GetReceiver())
	_, err = k.netw.SendMessage(&receiver, respMsg, false)
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
func (n *Network) SendPingMessage(tar, me *Contact, reqID int32) (*MessageHandler, error) {
	// TODO
	return nil, nil
}

func (n *Network) SendFindContactMessage(c *Contact) (*MessageHandler, error) {
	// TODO
	return nil, nil
}

func (n *Network) SendFindDataMessage(hash string) (*MessageHandler, error) {
	// TODO
	return nil, nil
}

func (n *Network) SendStoreMessage(data []byte) (*MessageHandler, error) {
	// TODO
	return nil, nil
}
