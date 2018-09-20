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
	Key       []byte              `protobuf:"bytes,5,opt,name=key,proto3" json:"key,omitempty"`
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

func (m *Message) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
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
	Id   int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Addr string `protobuf:"bytes,2,opt,name=addr" json:"addr,omitempty"`
}

func (m *Peer) Reset()                    { *m = Peer{} }
func (m *Peer) String() string            { return proto.CompactTextString(m) }
func (*Peer) ProtoMessage()               {}
func (*Peer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Peer) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
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
	// 340 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x51, 0x4d, 0x6f, 0xda, 0x40,
	0x10, 0xed, 0xda, 0x06, 0xec, 0x81, 0x5a, 0xd6, 0x5c, 0xba, 0x6a, 0x39, 0xac, 0xac, 0x1e, 0xac,
	0xb6, 0xe2, 0x40, 0xfb, 0x07, 0x90, 0x4c, 0x2b, 0x4b, 0x0d, 0xa0, 0x0d, 0xc9, 0x15, 0x2d, 0x78,
	0x14, 0x59, 0x60, 0xec, 0x78, 0x37, 0x91, 0xf8, 0xc3, 0xf9, 0x1d, 0xd1, 0x1a, 0x02, 0x21, 0xa7,
	0x7d, 0xf3, 0xde, 0xdb, 0xa7, 0xf9, 0x80, 0x70, 0xab, 0x72, 0x2a, 0x77, 0x85, 0x1a, 0xd5, 0x4d,
	0x65, 0x2a, 0x74, 0xea, 0x75, 0xfc, 0xe2, 0x40, 0xef, 0x86, 0xb4, 0x56, 0x0f, 0x84, 0x43, 0x08,
	0xca, 0x23, 0xcc, 0x52, 0xce, 0x04, 0x4b, 0x3a, 0xf2, 0x42, 0x58, 0xb5, 0xa1, 0xc7, 0x27, 0xd2,
	0x26, 0x4b, 0xb9, 0x73, 0x54, 0xcf, 0x04, 0xfe, 0x04, 0xcf, 0x1c, 0x6a, 0xe2, 0xae, 0x60, 0x49,
	0x38, 0xfe, 0x32, 0xaa, 0xd7, 0xa3, 0x53, 0xec, 0xdb, 0xbb, 0x3c, 0xd4, 0x24, 0x5b, 0x13, 0x7e,
	0x05, 0xbf, 0x21, 0x5d, 0x57, 0x7b, 0x4d, 0xdc, 0x13, 0x2c, 0xf1, 0xe5, 0xb9, 0xc6, 0x08, 0xdc,
	0x2d, 0x1d, 0x78, 0x47, 0xb0, 0x64, 0x20, 0x2d, 0x44, 0x01, 0x5d, 0x4d, 0xfb, 0x9c, 0x1a, 0xde,
	0x15, 0x2c, 0xe9, 0x8f, 0x7d, 0x1b, 0xbe, 0x20, 0x6a, 0xe4, 0x89, 0xc7, 0xef, 0x36, 0x6f, 0x43,
	0xc5, 0x33, 0x35, 0xbc, 0xf7, 0xc1, 0x73, 0x56, 0x70, 0x08, 0x5e, 0xae, 0x8c, 0xe2, 0xfe, 0xc5,
	0x91, 0x2a, 0xa3, 0x64, 0xcb, 0xe2, 0x37, 0x08, 0x34, 0xed, 0xcd, 0xca, 0x14, 0x25, 0xf1, 0x40,
	0xb0, 0xc4, 0x95, 0xbe, 0x25, 0x96, 0x45, 0x49, 0xf1, 0x04, 0xfa, 0xef, 0xa6, 0x40, 0x1f, 0xbc,
	0x45, 0x36, 0xfb, 0x17, 0x7d, 0xc2, 0xcf, 0x10, 0xfc, 0xcd, 0x66, 0xe9, 0x6a, 0x36, 0x4f, 0xa7,
	0x11, 0xc3, 0x10, 0xa0, 0x2d, 0xef, 0x27, 0xff, 0xef, 0xa6, 0x91, 0x83, 0x01, 0x74, 0x6e, 0x97,
	0x73, 0x39, 0x8d, 0xdc, 0xf8, 0x07, 0x78, 0xb6, 0x1f, 0x0c, 0xc1, 0x29, 0xf2, 0xd3, 0x76, 0x9d,
	0x22, 0x47, 0x04, 0x4f, 0xe5, 0x79, 0xd3, 0x6e, 0x34, 0x90, 0x2d, 0x8e, 0xff, 0x80, 0x67, 0x3b,
	0xc3, 0x5f, 0x30, 0xd8, 0xec, 0x2a, 0x4d, 0xda, 0xd8, 0xaf, 0x9a, 0x33, 0xe1, 0x5e, 0xcd, 0x76,
	0xa5, 0xae, 0xbb, 0xed, 0x55, 0x7f, 0xbf, 0x06, 0x00, 0x00, 0xff, 0xff, 0xb0, 0xa3, 0xcd, 0x57,
	0xe7, 0x01, 0x00, 0x00,
}