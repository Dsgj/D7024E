package pb

import (
	"time"
)

type MessageFactory struct {
	count int32
}

func NewMessageFactory() *MessageFactory {
	return &MessageFactory{count: 0}
}

func (msgFct *MessageFactory) new() int32 {
	current := msgFct.count
	msgFct.count++
	return current
}

func (msgFct *MessageFactory) NewMessage(t Message_MessageType, reqID int32,
	key string, sender, receiver *Peer, isResponse bool) *Message {
	time := time.Now().Unix()
	msgID := msgFct.new()
	msg := &Message{MessageID: msgID,
		Type:      t,
		RequestID: reqID,
		Key:       key,
		Sender:    sender,
		Receiver:  receiver,
		SentTime:  time,
		Response:  isResponse,
	}
	return msg
}

func (msgFct *MessageFactory) NewPingMessage(reqID int32,
	key string, sender, receiver *Peer, isResponse bool) *Message {
	return msgFct.NewMessage(Message_PING, reqID, key, sender, receiver, isResponse)
}

func (msgFct *MessageFactory) NewFindNodeMessage(reqID int32, key string,
	sender, receiver *Peer, isResponse bool) *Message {
	return msgFct.NewMessage(Message_FIND_NODE, reqID, key, sender, receiver, isResponse)
}

func (msg *Message) AddPeerData(peers []*Peer) *Message {
	msg.Data = &Data{ClosestPeers: peers}
	return msg
}
