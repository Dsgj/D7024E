// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kademlia.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	kademlia.proto

It has these top-level messages:
	Message
	Peer
	Record
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
	Id       string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Addr     string `protobuf:"bytes,2,opt,name=addr" json:"addr,omitempty"`
	Distance string `protobuf:"bytes,3,opt,name=distance" json:"distance,omitempty"`
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

func (m *Peer) GetDistance() string {
	if m != nil {
		return m.Distance
	}
	return ""
}

type Record struct {
	Key         []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value       []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	NewPublish  bool   `protobuf:"varint,3,opt,name=newPublish" json:"newPublish,omitempty"`
	Publisher   *Peer  `protobuf:"bytes,4,opt,name=publisher" json:"publisher,omitempty"`
	PublishedAt int64  `protobuf:"varint,5,opt,name=publishedAt" json:"publishedAt,omitempty"`
}

func (m *Record) Reset()                    { *m = Record{} }
func (m *Record) String() string            { return proto.CompactTextString(m) }
func (*Record) ProtoMessage()               {}
func (*Record) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Record) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Record) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Record) GetNewPublish() bool {
	if m != nil {
		return m.NewPublish
	}
	return false
}

func (m *Record) GetPublisher() *Peer {
	if m != nil {
		return m.Publisher
	}
	return nil
}

func (m *Record) GetPublishedAt() int64 {
	if m != nil {
		return m.PublishedAt
	}
	return 0
}

type Data struct {
	// FIND_NODE
	ClosestPeers []*Peer `protobuf:"bytes,1,rep,name=closestPeers" json:"closestPeers,omitempty"`
	// FIND_VALUE
	Record *Record `protobuf:"bytes,2,opt,name=record" json:"record,omitempty"`
}

func (m *Data) Reset()                    { *m = Data{} }
func (m *Data) String() string            { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()               {}
func (*Data) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Data) GetClosestPeers() []*Peer {
	if m != nil {
		return m.ClosestPeers
	}
	return nil
}

func (m *Data) GetRecord() *Record {
	if m != nil {
		return m.Record
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "pb.Message")
	proto.RegisterType((*Peer)(nil), "pb.Peer")
	proto.RegisterType((*Record)(nil), "pb.Record")
	proto.RegisterType((*Data)(nil), "pb.Data")
	proto.RegisterEnum("pb.Message_MessageType", Message_MessageType_name, Message_MessageType_value)
}

func init() { proto.RegisterFile("kademlia.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 436 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x52, 0xdd, 0x8a, 0xd3, 0x40,
	0x14, 0x76, 0x92, 0xb4, 0x9b, 0x39, 0xad, 0xa5, 0x1c, 0x04, 0x07, 0x5d, 0x24, 0x04, 0x91, 0x80,
	0xd2, 0x8b, 0xfa, 0x04, 0x85, 0x76, 0xa5, 0xa0, 0xdd, 0x32, 0x56, 0xf1, 0x6e, 0x99, 0x76, 0x0e,
	0x1a, 0xb6, 0x6d, 0xe2, 0xcc, 0x74, 0xa5, 0xaf, 0xe2, 0x03, 0xfa, 0x1c, 0x32, 0x93, 0xfe, 0xed,
	0x5e, 0xe5, 0x7c, 0x3f, 0x7c, 0xcc, 0xf9, 0x4e, 0xa0, 0x77, 0xaf, 0x34, 0x6d, 0xd6, 0xa5, 0x1a,
	0xd4, 0xa6, 0x72, 0x15, 0x46, 0xf5, 0x32, 0xff, 0x17, 0xc1, 0xd5, 0x17, 0xb2, 0x56, 0xfd, 0x24,
	0xbc, 0x06, 0xbe, 0x69, 0xc6, 0xe9, 0x58, 0xb0, 0x8c, 0x15, 0x2d, 0x79, 0x26, 0xbc, 0x6a, 0xe8,
	0xf7, 0x8e, 0xac, 0x9b, 0x8e, 0x45, 0xd4, 0xa8, 0x27, 0x02, 0xdf, 0x43, 0xe2, 0xf6, 0x35, 0x89,
	0x38, 0x63, 0x45, 0x6f, 0xf8, 0x72, 0x50, 0x2f, 0x07, 0x87, 0xd8, 0xe3, 0x77, 0xb1, 0xaf, 0x49,
	0x06, 0x13, 0xbe, 0x82, 0xd4, 0x90, 0xad, 0xab, 0xad, 0x25, 0x91, 0x64, 0xac, 0x48, 0xe5, 0x09,
	0x63, 0x1f, 0xe2, 0x7b, 0xda, 0x8b, 0x56, 0xc6, 0x0a, 0x2e, 0xfd, 0x88, 0x19, 0xb4, 0x2d, 0x6d,
	0x35, 0x19, 0xd1, 0xce, 0x58, 0xd1, 0x19, 0xa6, 0x3e, 0x7c, 0x4e, 0x64, 0xe4, 0x81, 0xc7, 0xb7,
	0x3e, 0x6f, 0x45, 0xe5, 0x03, 0x19, 0x71, 0xf5, 0xc4, 0x73, 0x52, 0xf0, 0x1a, 0x12, 0xad, 0x9c,
	0x12, 0xe9, 0xd9, 0x31, 0x56, 0x4e, 0xc9, 0xc0, 0xe2, 0x6b, 0xe0, 0x96, 0xb6, 0xee, 0xce, 0x95,
	0x1b, 0x12, 0x3c, 0x63, 0x45, 0x2c, 0x53, 0x4f, 0x2c, 0xca, 0x0d, 0xe5, 0x23, 0xe8, 0x5c, 0x6c,
	0x81, 0x29, 0x24, 0xf3, 0xe9, 0xec, 0x53, 0xff, 0x19, 0x3e, 0x07, 0x7e, 0x33, 0x9d, 0x8d, 0xef,
	0x66, 0xb7, 0xe3, 0x49, 0x9f, 0x61, 0x0f, 0x20, 0xc0, 0xef, 0xa3, 0xcf, 0xdf, 0x26, 0xfd, 0x08,
	0x39, 0xb4, 0xbe, 0x2e, 0x6e, 0xe5, 0xa4, 0x1f, 0xe7, 0x37, 0x90, 0xf8, 0xf7, 0x60, 0x0f, 0xa2,
	0x52, 0x87, 0x76, 0xb9, 0x8c, 0x4a, 0x8d, 0x08, 0x89, 0xd2, 0xda, 0x84, 0x46, 0xb9, 0x0c, 0xb3,
	0xef, 0x47, 0x97, 0xd6, 0xa9, 0xed, 0xaa, 0x29, 0x94, 0xcb, 0x13, 0xce, 0xff, 0x32, 0x68, 0x4b,
	0x5a, 0x55, 0x46, 0x1f, 0xab, 0xf2, 0x59, 0xdd, 0xa6, 0xaa, 0x17, 0xd0, 0x7a, 0x50, 0xeb, 0x1d,
	0x85, 0xb4, 0xae, 0x6c, 0x00, 0xbe, 0x01, 0xd8, 0xd2, 0x9f, 0xf9, 0x6e, 0xb9, 0x2e, 0xed, 0xaf,
	0x10, 0x98, 0xca, 0x0b, 0x06, 0xdf, 0x01, 0xaf, 0x9b, 0x91, 0x4c, 0xb8, 0xc7, 0x65, 0x7f, 0x67,
	0x09, 0x33, 0xe8, 0x1c, 0x81, 0x1e, 0xb9, 0x70, 0xa2, 0x58, 0x5e, 0x52, 0xf9, 0x0f, 0x48, 0x7c,
	0xa5, 0xf8, 0x01, 0xba, 0xab, 0x75, 0x65, 0xc9, 0x3a, 0x9f, 0x61, 0x05, 0xcb, 0xe2, 0x47, 0xa1,
	0x8f, 0x54, 0xcc, 0xa1, 0x6d, 0xc2, 0x46, 0xe1, 0xd9, 0x9d, 0x21, 0x78, 0x5f, 0xb3, 0xa3, 0x3c,
	0x28, 0xcb, 0x76, 0xf8, 0x65, 0x3f, 0xfe, 0x0f, 0x00, 0x00, 0xff, 0xff, 0xf4, 0xe5, 0x6c, 0xdd,
	0xc4, 0x02, 0x00, 0x00,
}
