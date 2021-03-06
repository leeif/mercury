// Code generated by protoc-gen-go. DO NOT EDIT.
// source: index.proto

package index

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type RegisterMemberRequest struct {
	Mid                  string   `protobuf:"bytes,1,opt,name=mid,proto3" json:"mid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterMemberRequest) Reset()         { *m = RegisterMemberRequest{} }
func (m *RegisterMemberRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterMemberRequest) ProtoMessage()    {}
func (*RegisterMemberRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f750e0f7889345b5, []int{0}
}

func (m *RegisterMemberRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterMemberRequest.Unmarshal(m, b)
}
func (m *RegisterMemberRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterMemberRequest.Marshal(b, m, deterministic)
}
func (m *RegisterMemberRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterMemberRequest.Merge(m, src)
}
func (m *RegisterMemberRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterMemberRequest.Size(m)
}
func (m *RegisterMemberRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterMemberRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterMemberRequest proto.InternalMessageInfo

func (m *RegisterMemberRequest) GetMid() string {
	if m != nil {
		return m.Mid
	}
	return ""
}

type RegisterMemberReply struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterMemberReply) Reset()         { *m = RegisterMemberReply{} }
func (m *RegisterMemberReply) String() string { return proto.CompactTextString(m) }
func (*RegisterMemberReply) ProtoMessage()    {}
func (*RegisterMemberReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_f750e0f7889345b5, []int{1}
}

func (m *RegisterMemberReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterMemberReply.Unmarshal(m, b)
}
func (m *RegisterMemberReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterMemberReply.Marshal(b, m, deterministic)
}
func (m *RegisterMemberReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterMemberReply.Merge(m, src)
}
func (m *RegisterMemberReply) XXX_Size() int {
	return xxx_messageInfo_RegisterMemberReply.Size(m)
}
func (m *RegisterMemberReply) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterMemberReply.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterMemberReply proto.InternalMessageInfo

func (m *RegisterMemberReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type GetMemberRequest struct {
	Mid                  string   `protobuf:"bytes,1,opt,name=mid,proto3" json:"mid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetMemberRequest) Reset()         { *m = GetMemberRequest{} }
func (m *GetMemberRequest) String() string { return proto.CompactTextString(m) }
func (*GetMemberRequest) ProtoMessage()    {}
func (*GetMemberRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f750e0f7889345b5, []int{2}
}

func (m *GetMemberRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetMemberRequest.Unmarshal(m, b)
}
func (m *GetMemberRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetMemberRequest.Marshal(b, m, deterministic)
}
func (m *GetMemberRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetMemberRequest.Merge(m, src)
}
func (m *GetMemberRequest) XXX_Size() int {
	return xxx_messageInfo_GetMemberRequest.Size(m)
}
func (m *GetMemberRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetMemberRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetMemberRequest proto.InternalMessageInfo

func (m *GetMemberRequest) GetMid() string {
	if m != nil {
		return m.Mid
	}
	return ""
}

type GetMemberReply struct {
	Mid                  string   `protobuf:"bytes,1,opt,name=mid,proto3" json:"mid,omitempty"`
	Ip                   string   `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetMemberReply) Reset()         { *m = GetMemberReply{} }
func (m *GetMemberReply) String() string { return proto.CompactTextString(m) }
func (*GetMemberReply) ProtoMessage()    {}
func (*GetMemberReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_f750e0f7889345b5, []int{3}
}

func (m *GetMemberReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetMemberReply.Unmarshal(m, b)
}
func (m *GetMemberReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetMemberReply.Marshal(b, m, deterministic)
}
func (m *GetMemberReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetMemberReply.Merge(m, src)
}
func (m *GetMemberReply) XXX_Size() int {
	return xxx_messageInfo_GetMemberReply.Size(m)
}
func (m *GetMemberReply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetMemberReply.DiscardUnknown(m)
}

var xxx_messageInfo_GetMemberReply proto.InternalMessageInfo

func (m *GetMemberReply) GetMid() string {
	if m != nil {
		return m.Mid
	}
	return ""
}

func (m *GetMemberReply) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func init() {
	proto.RegisterType((*RegisterMemberRequest)(nil), "index.RegisterMemberRequest")
	proto.RegisterType((*RegisterMemberReply)(nil), "index.RegisterMemberReply")
	proto.RegisterType((*GetMemberRequest)(nil), "index.GetMemberRequest")
	proto.RegisterType((*GetMemberReply)(nil), "index.GetMemberReply")
}

func init() { proto.RegisterFile("index.proto", fileDescriptor_f750e0f7889345b5) }

var fileDescriptor_f750e0f7889345b5 = []byte{
	// 197 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0xcc, 0x4b, 0x49,
	0xad, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0x94, 0x34, 0xb9, 0x44, 0x83,
	0x52, 0xd3, 0x33, 0x8b, 0x4b, 0x52, 0x8b, 0x7c, 0x53, 0x73, 0x93, 0x52, 0x8b, 0x82, 0x52, 0x0b,
	0x4b, 0x53, 0x8b, 0x4b, 0x84, 0x04, 0xb8, 0x98, 0x73, 0x33, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35,
	0x38, 0x83, 0x40, 0x4c, 0x25, 0x7d, 0x2e, 0x61, 0x74, 0xa5, 0x05, 0x39, 0x95, 0x42, 0x12, 0x5c,
	0xec, 0xb9, 0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9, 0x50, 0xc5, 0x30, 0xae, 0x92, 0x0a, 0x97, 0x80,
	0x7b, 0x6a, 0x09, 0x21, 0x63, 0x8d, 0xb8, 0xf8, 0x90, 0x54, 0x81, 0x4c, 0xc4, 0x50, 0x23, 0xc4,
	0xc7, 0xc5, 0x94, 0x59, 0x20, 0xc1, 0x04, 0x16, 0x60, 0xca, 0x2c, 0x30, 0x9a, 0xc1, 0xc8, 0xc5,
	0xe9, 0x09, 0x72, 0xbf, 0x5f, 0x7e, 0x4a, 0xaa, 0x90, 0x0f, 0x17, 0x1f, 0xaa, 0xc3, 0x84, 0x64,
	0xf4, 0x20, 0x5e, 0xc5, 0xea, 0x35, 0x29, 0x29, 0x1c, 0xb2, 0x05, 0x39, 0x95, 0x4a, 0x0c, 0x42,
	0xb6, 0x5c, 0x9c, 0x70, 0xf7, 0x08, 0x89, 0x43, 0x95, 0xa2, 0xfb, 0x43, 0x4a, 0x14, 0x53, 0x02,
	0xac, 0x3d, 0x89, 0x0d, 0x1c, 0xbc, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x68, 0xff, 0x38,
	0xe7, 0x6d, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// IndexNodeClient is the client API for IndexNode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IndexNodeClient interface {
	RegisterMember(ctx context.Context, in *RegisterMemberRequest, opts ...grpc.CallOption) (*RegisterMemberReply, error)
	GetMember(ctx context.Context, in *GetMemberRequest, opts ...grpc.CallOption) (*GetMemberReply, error)
}

type indexNodeClient struct {
	cc *grpc.ClientConn
}

func NewIndexNodeClient(cc *grpc.ClientConn) IndexNodeClient {
	return &indexNodeClient{cc}
}

func (c *indexNodeClient) RegisterMember(ctx context.Context, in *RegisterMemberRequest, opts ...grpc.CallOption) (*RegisterMemberReply, error) {
	out := new(RegisterMemberReply)
	err := c.cc.Invoke(ctx, "/index.IndexNode/RegisterMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexNodeClient) GetMember(ctx context.Context, in *GetMemberRequest, opts ...grpc.CallOption) (*GetMemberReply, error) {
	out := new(GetMemberReply)
	err := c.cc.Invoke(ctx, "/index.IndexNode/GetMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IndexNodeServer is the server API for IndexNode service.
type IndexNodeServer interface {
	RegisterMember(context.Context, *RegisterMemberRequest) (*RegisterMemberReply, error)
	GetMember(context.Context, *GetMemberRequest) (*GetMemberReply, error)
}

// UnimplementedIndexNodeServer can be embedded to have forward compatible implementations.
type UnimplementedIndexNodeServer struct {
}

func (*UnimplementedIndexNodeServer) RegisterMember(ctx context.Context, req *RegisterMemberRequest) (*RegisterMemberReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterMember not implemented")
}
func (*UnimplementedIndexNodeServer) GetMember(ctx context.Context, req *GetMemberRequest) (*GetMemberReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMember not implemented")
}

func RegisterIndexNodeServer(s *grpc.Server, srv IndexNodeServer) {
	s.RegisterService(&_IndexNode_serviceDesc, srv)
}

func _IndexNode_RegisterMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexNodeServer).RegisterMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/index.IndexNode/RegisterMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexNodeServer).RegisterMember(ctx, req.(*RegisterMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndexNode_GetMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexNodeServer).GetMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/index.IndexNode/GetMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexNodeServer).GetMember(ctx, req.(*GetMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _IndexNode_serviceDesc = grpc.ServiceDesc{
	ServiceName: "index.IndexNode",
	HandlerType: (*IndexNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterMember",
			Handler:    _IndexNode_RegisterMember_Handler,
		},
		{
			MethodName: "GetMember",
			Handler:    _IndexNode_GetMember_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "index.proto",
}
