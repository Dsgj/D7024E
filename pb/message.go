package pb


func NewMessage(t Message_MessageType, reqID int32,
                key []byte, sender, receiver *Peer) *Message {
    msg := &Message{Type:       t,
                    RequestID:  reqID,
                    Key:        key,
                    Sender:     sender,
                    Receiver:   receiver,
    }
    return msg
}
