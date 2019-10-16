// Code generated by protoc-gen-go. DO NOT EDIT.
// source: addressbook.proto

package gen

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Person_PhoneType int32

const (
	Person_MOBILE Person_PhoneType = 0
	Person_HOME   Person_PhoneType = 1
	Person_WORK   Person_PhoneType = 2
)

var Person_PhoneType_name = map[int32]string{
	0: "MOBILE",
	1: "HOME",
	2: "WORK",
}

var Person_PhoneType_value = map[string]int32{
	"MOBILE": 0,
	"HOME":   1,
	"WORK":   2,
}

func (x Person_PhoneType) String() string {
	return proto.EnumName(Person_PhoneType_name, int32(x))
}

func (Person_PhoneType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_1eb1a68c9dd6d429, []int{0, 0}
}

// [START messages]
type Person struct {
	Name                 string                `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Id                   int32                 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Email                string                `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Phones               []*Person_PhoneNumber `protobuf:"bytes,4,rep,name=phones,proto3" json:"phones,omitempty"`
	LastUpdated          *timestamp.Timestamp  `protobuf:"bytes,5,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Person) Reset()         { *m = Person{} }
func (m *Person) String() string { return proto.CompactTextString(m) }
func (*Person) ProtoMessage()    {}
func (*Person) Descriptor() ([]byte, []int) {
	return fileDescriptor_1eb1a68c9dd6d429, []int{0}
}

func (m *Person) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Person.Unmarshal(m, b)
}
func (m *Person) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Person.Marshal(b, m, deterministic)
}
func (m *Person) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Person.Merge(m, src)
}
func (m *Person) XXX_Size() int {
	return xxx_messageInfo_Person.Size(m)
}
func (m *Person) XXX_DiscardUnknown() {
	xxx_messageInfo_Person.DiscardUnknown(m)
}

var xxx_messageInfo_Person proto.InternalMessageInfo

func (m *Person) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Person) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Person) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Person) GetPhones() []*Person_PhoneNumber {
	if m != nil {
		return m.Phones
	}
	return nil
}

func (m *Person) GetLastUpdated() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdated
	}
	return nil
}

type Person_PhoneNumber struct {
	Number               string           `protobuf:"bytes,1,opt,name=number,proto3" json:"number,omitempty"`
	Type                 Person_PhoneType `protobuf:"varint,2,opt,name=type,proto3,enum=gen.Person_PhoneType" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Person_PhoneNumber) Reset()         { *m = Person_PhoneNumber{} }
func (m *Person_PhoneNumber) String() string { return proto.CompactTextString(m) }
func (*Person_PhoneNumber) ProtoMessage()    {}
func (*Person_PhoneNumber) Descriptor() ([]byte, []int) {
	return fileDescriptor_1eb1a68c9dd6d429, []int{0, 0}
}

func (m *Person_PhoneNumber) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Person_PhoneNumber.Unmarshal(m, b)
}
func (m *Person_PhoneNumber) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Person_PhoneNumber.Marshal(b, m, deterministic)
}
func (m *Person_PhoneNumber) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Person_PhoneNumber.Merge(m, src)
}
func (m *Person_PhoneNumber) XXX_Size() int {
	return xxx_messageInfo_Person_PhoneNumber.Size(m)
}
func (m *Person_PhoneNumber) XXX_DiscardUnknown() {
	xxx_messageInfo_Person_PhoneNumber.DiscardUnknown(m)
}

var xxx_messageInfo_Person_PhoneNumber proto.InternalMessageInfo

func (m *Person_PhoneNumber) GetNumber() string {
	if m != nil {
		return m.Number
	}
	return ""
}

func (m *Person_PhoneNumber) GetType() Person_PhoneType {
	if m != nil {
		return m.Type
	}
	return Person_MOBILE
}

// Our address book file is just one of these.
type AddressBook struct {
	People               []*Person `protobuf:"bytes,1,rep,name=people,proto3" json:"people,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *AddressBook) Reset()         { *m = AddressBook{} }
func (m *AddressBook) String() string { return proto.CompactTextString(m) }
func (*AddressBook) ProtoMessage()    {}
func (*AddressBook) Descriptor() ([]byte, []int) {
	return fileDescriptor_1eb1a68c9dd6d429, []int{1}
}

func (m *AddressBook) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddressBook.Unmarshal(m, b)
}
func (m *AddressBook) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddressBook.Marshal(b, m, deterministic)
}
func (m *AddressBook) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddressBook.Merge(m, src)
}
func (m *AddressBook) XXX_Size() int {
	return xxx_messageInfo_AddressBook.Size(m)
}
func (m *AddressBook) XXX_DiscardUnknown() {
	xxx_messageInfo_AddressBook.DiscardUnknown(m)
}

var xxx_messageInfo_AddressBook proto.InternalMessageInfo

func (m *AddressBook) GetPeople() []*Person {
	if m != nil {
		return m.People
	}
	return nil
}

func init() {
	proto.RegisterEnum("gen.Person_PhoneType", Person_PhoneType_name, Person_PhoneType_value)
	proto.RegisterType((*Person)(nil), "gen.Person")
	proto.RegisterType((*Person_PhoneNumber)(nil), "gen.Person.PhoneNumber")
	proto.RegisterType((*AddressBook)(nil), "gen.AddressBook")
}

func init() { proto.RegisterFile("addressbook.proto", fileDescriptor_1eb1a68c9dd6d429) }

var fileDescriptor_1eb1a68c9dd6d429 = []byte{
	// 351 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0x4f, 0x6b, 0xa3, 0x40,
	0x18, 0xc6, 0x57, 0x63, 0x64, 0xf3, 0xba, 0x84, 0x64, 0xc8, 0xee, 0x4a, 0x2e, 0x2b, 0xd9, 0x3d,
	0xb8, 0x2c, 0x4c, 0xc0, 0x3d, 0xf7, 0x50, 0x21, 0xb4, 0xa5, 0x4d, 0x23, 0x92, 0xd2, 0x63, 0x19,
	0xeb, 0x5b, 0x2b, 0x51, 0x67, 0x70, 0x46, 0x68, 0xbe, 0x52, 0x6f, 0xfd, 0x86, 0x45, 0xc7, 0x14,
	0xa1, 0xb7, 0xf7, 0xcf, 0x6f, 0x9e, 0x79, 0xde, 0x07, 0xe6, 0x2c, 0x4d, 0x6b, 0x94, 0x32, 0xe1,
	0xfc, 0x40, 0x45, 0xcd, 0x15, 0x27, 0xa3, 0x0c, 0xab, 0xe5, 0xaf, 0x8c, 0xf3, 0xac, 0xc0, 0x75,
	0x37, 0x4a, 0x9a, 0xa7, 0xb5, 0xca, 0x4b, 0x94, 0x8a, 0x95, 0x42, 0x53, 0xab, 0x37, 0x13, 0xec,
	0x08, 0x6b, 0xc9, 0x2b, 0x42, 0xc0, 0xaa, 0x58, 0x89, 0xae, 0xe1, 0x19, 0xfe, 0x24, 0xee, 0x6a,
	0x32, 0x05, 0x33, 0x4f, 0x5d, 0xd3, 0x33, 0xfc, 0x71, 0x6c, 0xe6, 0x29, 0x59, 0xc0, 0x18, 0x4b,
	0x96, 0x17, 0xee, 0xa8, 0x83, 0x74, 0x43, 0xd6, 0x60, 0x8b, 0x67, 0x5e, 0xa1, 0x74, 0x2d, 0x6f,
	0xe4, 0x3b, 0xc1, 0x4f, 0x9a, 0x61, 0x45, 0xb5, 0x2c, 0x8d, 0xda, 0xcd, 0x6d, 0x53, 0x26, 0x58,
	0xc7, 0x3d, 0x46, 0xce, 0xe0, 0x5b, 0xc1, 0xa4, 0x7a, 0x68, 0x44, 0xca, 0x14, 0xa6, 0xee, 0xd8,
	0x33, 0x7c, 0x27, 0x58, 0x52, 0xed, 0x96, 0x9e, 0xdc, 0xd2, 0xfd, 0xc9, 0x6d, 0xec, 0xb4, 0xfc,
	0x9d, 0xc6, 0x97, 0x11, 0x38, 0x03, 0x55, 0xf2, 0x03, 0xec, 0xaa, 0xab, 0x7a, 0xeb, 0x7d, 0x47,
	0xfe, 0x82, 0xa5, 0x8e, 0x02, 0x3b, 0xfb, 0xd3, 0xe0, 0xfb, 0x27, 0x53, 0xfb, 0xa3, 0xc0, 0xb8,
	0x43, 0x56, 0xff, 0x60, 0xf2, 0x31, 0x22, 0x00, 0xf6, 0x76, 0x17, 0x5e, 0xdd, 0x6c, 0x66, 0x5f,
	0xc8, 0x57, 0xb0, 0x2e, 0x77, 0xdb, 0xcd, 0xcc, 0x68, 0xab, 0xfb, 0x5d, 0x7c, 0x3d, 0x33, 0x57,
	0x01, 0x38, 0xe7, 0x3a, 0xee, 0x90, 0xf3, 0x03, 0xf9, 0x0d, 0xb6, 0x40, 0x2e, 0x8a, 0x36, 0xb9,
	0xf6, 0x7a, 0x67, 0xf0, 0x51, 0xdc, 0xaf, 0xc2, 0x08, 0x16, 0x8f, 0xbc, 0xa4, 0xf8, 0xc2, 0x4a,
	0x51, 0x20, 0x55, 0x8d, 0xe2, 0x75, 0xce, 0x8a, 0x70, 0x3e, 0x50, 0x8a, 0xda, 0xb3, 0xe5, 0xab,
	0xf9, 0xe7, 0x42, 0xc7, 0x10, 0x9d, 0x62, 0xd8, 0xe8, 0x57, 0x92, 0x0e, 0xe0, 0xc4, 0xee, 0x52,
	0xfa, 0xff, 0x1e, 0x00, 0x00, 0xff, 0xff, 0x5f, 0xa9, 0xb5, 0x04, 0xfb, 0x01, 0x00, 0x00,
}
