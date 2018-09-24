// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kademlia.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	kademlia.proto

It has these top-level messages:
	Message
	Peer
	Data
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Message_MessageType int32

const (
	Message_PING       Message_MessageType = 0
	Message_FIND_NODE  Message_MessageType = 1
	Message_FIND_VALUE Message_MessageType = 2
	Message_STORE      Message_MessageType = 3
)

var Message_MessageType_name = map[int32]string{
	0: "PING",
	1: "FIND_NODE",
	2: "FIND_VALUE",
	3: "STORE",
}
var Message_MessageType_value = map[string]int32{
	"PING":       0,
	"FIND_NODE":  1,
	"FIND_VALUE": 2,
	"STORE":      3,
}

func (x Message_MessageType) String() string {
	return proto.EnumName(Message_MessageType_name, int32(x))
}
func (Message_MessageType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Message struct {
	MessageID int32               `protobuf:"varint,1,opt,name=messageID" json:"messageID,omitempty"`
	RequestID int32               `protobuf:"varint,2,opt,name=requestID" json:"requestID,omitempty"`
	Type      Message_MessageType `protobuf:"varint,3,opt,name=type,enum=pb.Message_MessageType" json:"type,omitempty"`
	Response  bool                `protobuf:"varint,4,opt,name=response" json:"response,omitempty"`
	Key       string              `protobuf:"bytes,5,opt,name=key" json:"key,omitempty"`
	Sender    *Peer               `protobuf:"bytes,6,opt,name=sender" json:"sender,omitempty"`
	Receiver  *Peer               `protobuf:"bytes,7,opt,name=receiver" json:"receiver,omitempty"`
	Data      *Data               `protobuf:"bytes,8,opt,name=data" json:"data,omitempty"`
	SentTime  int64               `protobuf:"varint,9,opt,name=sent_time,json=sentTime" json:"sent_time,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetMessageID() int32 {
	if m != nil {
		return m.MessageID
	}
	return 0
}

func (m *Message) GetRequestID() int32 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *Message) GetType() Message_MessageType {
	if m != nil {
		return m.Type
	}
	return Message_PING
}

func (m *Message) GetResponse() bool {
	if m != nil {
		return m.Response
	}
	return false
}

func (m *Message) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Message) GetSender() *Peer {
	if m != nil {
		return m.Sender
	}
	return nil
}

func (m *Message) GetReceiver() *Peer {
	if m != nil {
		return m.Receiver
	}
	return nil
}

func (m *Message) GetData() *Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Message) GetSentTime() int64 {
	if m != nil {
		return m.SentTime
	}
	return 0
}

type Peer struct {
	Id   string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Addr string `protobuf:"bytes,2,opt,name=addr" json:"addr,omitempty"`
}

func (m *Peer) Reset()                    { *m = Peer{} }
func (m *Peer) String() string            { return proto.CompactTextString(m) }
func (*Peer) ProtoMessage()               {}
func (*Peer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Peer) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Peer) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type Data struct {
	// FIND_NODE
	ClosestPeers []*Peer `protobuf:"bytes,1,rep,name=closestPeers" json:"closestPeers,omitempty"`
}

func (m *Data) Reset()                    { *m = Data{} }
func (m *Data) String() string            { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()               {}
func (*Data) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Data) GetClosestPeers() []*Peer {
	if m != nil {
		return m.ClosestPeers
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "pb.Message")
	proto.RegisterType((*Peer)(nil), "pb.Peer")
	proto.RegisterType((*Data)(nil), "pb.Data")
	proto.RegisterEnum("pb.Message_MessageType", Message_MessageType_name, Message_MessageType_value)
}

func init() { proto.RegisterFile("kademlia.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 341 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x51, 0x4d, 0x8f, 0xda, 0x30,
	0x10, 0xad, 0x93, 0x00, 0xf1, 0xd0, 0x46, 0xd1, 0x5c, 0x6a, 0xb5, 0x1c, 0xac, 0xa8, 0x87, 0xa8,
	0xad, 0x38, 0xd0, 0xfe, 0x01, 0xa4, 0xd0, 0x2a, 0x52, 0x0b, 0xc8, 0xa5, 0xbd, 0x22, 0x43, 0x46,
	0x55, 0x04, 0x21, 0xd9, 0xd8, 0xbb, 0x12, 0x7f, 0x78, 0x7f, 0xc7, 0xca, 0x81, 0x85, 0x65, 0x4f,
	0x7e, 0xf3, 0xde, 0xf3, 0xd3, 0x7c, 0x40, 0xb4, 0xd3, 0x05, 0x55, 0xfb, 0x52, 0x8f, 0x9b, 0xb6,
	0xb6, 0x35, 0x7a, 0xcd, 0x26, 0x79, 0xf4, 0x60, 0xf0, 0x9b, 0x8c, 0xd1, 0xff, 0x09, 0x47, 0xc0,
	0xab, 0x13, 0xcc, 0x33, 0xc1, 0x24, 0x4b, 0x7b, 0xea, 0x4a, 0x38, 0xb5, 0xa5, 0xbb, 0x7b, 0x32,
	0x36, 0xcf, 0x84, 0x77, 0x52, 0x2f, 0x04, 0x7e, 0x81, 0xc0, 0x1e, 0x1b, 0x12, 0xbe, 0x64, 0x69,
	0x34, 0x79, 0x3f, 0x6e, 0x36, 0xe3, 0x73, 0xec, 0xf3, 0xbb, 0x3a, 0x36, 0xa4, 0x3a, 0x13, 0x7e,
	0x80, 0xb0, 0x25, 0xd3, 0xd4, 0x07, 0x43, 0x22, 0x90, 0x2c, 0x0d, 0xd5, 0xa5, 0xc6, 0x18, 0xfc,
	0x1d, 0x1d, 0x45, 0x4f, 0xb2, 0x94, 0x2b, 0x07, 0x51, 0x42, 0xdf, 0xd0, 0xa1, 0xa0, 0x56, 0xf4,
	0x25, 0x4b, 0x87, 0x93, 0xd0, 0x85, 0x2f, 0x89, 0x5a, 0x75, 0xe6, 0xf1, 0x93, 0xcb, 0xdb, 0x52,
	0xf9, 0x40, 0xad, 0x18, 0xbc, 0xf2, 0x5c, 0x14, 0x1c, 0x41, 0x50, 0x68, 0xab, 0x45, 0x78, 0x75,
	0x64, 0xda, 0x6a, 0xd5, 0xb1, 0xf8, 0x11, 0xb8, 0xa1, 0x83, 0x5d, 0xdb, 0xb2, 0x22, 0xc1, 0x25,
	0x4b, 0x7d, 0x15, 0x3a, 0x62, 0x55, 0x56, 0x94, 0x4c, 0x61, 0xf8, 0x62, 0x0a, 0x0c, 0x21, 0x58,
	0xe6, 0xf3, 0x9f, 0xf1, 0x1b, 0x7c, 0x07, 0xfc, 0x47, 0x3e, 0xcf, 0xd6, 0xf3, 0x45, 0x36, 0x8b,
	0x19, 0x46, 0x00, 0x5d, 0xf9, 0x6f, 0xfa, 0xeb, 0xef, 0x2c, 0xf6, 0x90, 0x43, 0xef, 0xcf, 0x6a,
	0xa1, 0x66, 0xb1, 0x9f, 0x7c, 0x86, 0xc0, 0xf5, 0x83, 0x11, 0x78, 0x65, 0xd1, 0x6d, 0x97, 0x2b,
	0xaf, 0x2c, 0x10, 0x21, 0xd0, 0x45, 0xd1, 0x76, 0x1b, 0xe5, 0xaa, 0xc3, 0xc9, 0x77, 0x08, 0x5c,
	0x67, 0xf8, 0x15, 0xde, 0x6e, 0xf7, 0xb5, 0x21, 0x63, 0xdd, 0x57, 0x23, 0x98, 0xf4, 0x6f, 0x66,
	0xbb, 0x51, 0x37, 0xfd, 0xee, 0xaa, 0xdf, 0x9e, 0x02, 0x00, 0x00, 0xff, 0xff, 0x5f, 0x0c, 0xfd,
	0x7e, 0xe7, 0x01, 0x00, 0x00,
}
