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
	msgCounter  int32
}

func NewNetwork(port, addr string) *Network {
	netw := &Network{port: port,
		addr:        addr,
		msgHandlers: make(map[int32]*MessageHandler),
		msgCounter:  0}
	return netw
}

type MessageHandler struct {
	id int32
	ch chan (*pb.Message)
}

func (mh *MessageHandler) awaitMessage(returnChan chan (*pb.Message)) {
	defer close(mh.ch)
	select {
	case result := <-mh.ch:
		returnChan <- result
		return
	case <-time.After(30 * time.Second):
		//handle timed out message
		return
	}
}

func Listen(k *Kademlia, ip string, port string) error {
	network := k.netw
	serverAddr, err := net.ResolveUDPAddr("udp", ":"+network.port)
	if err != nil {
		return err
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return err
	}
	defer serverConn.Close()

	buf := make([]byte, 1024)

	log.Println("Listening on port " + network.port)
	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		msg := &pb.Message{}
		err = proto.Unmarshal(buf[0:n], msg)
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
		} else {
			//handler := k.getTypeHandler(msg.GetType())
			//HANDLE IT
		}
	}
}

func (n *Network) SendMessage(c *Contact,
	msg *pb.Message,
	wantResponse bool) (*MessageHandler, error) {
	// Should probably create and manage connections in a better way
	remoteAddr, err := net.ResolveUDPAddr("udp", c.Address+":"+n.port)
	if err != nil {
		return nil, err
	}
	localAddr, err := net.ResolveUDPAddr("udp", n.addr+":"+n.port)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}

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

func (n *Network) SendPingMessage(tar, me *Contact, reqID int32) (*MessageHandler, error) {
	//create ping message
	msg := pb.NewMessage(pb.Message_PING, reqID, nil,
		ContactInfoToPeer(me), ContactInfoToPeer(tar))
	msg.MessageID = n.msgCounter
	n.msgCounter++
	//call SendMessage
	return n.SendMessage(tar, msg, true)
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
