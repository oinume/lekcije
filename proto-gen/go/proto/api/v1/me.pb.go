// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/api/v1/me.proto

/*
Package api_v1 is a generated protocol buffer package.

It is generated from these files:
	proto/api/v1/me.proto

It has these top-level messages:
	GetMeEmailRequest
	GetMeEmailResponse
	UpdateMeEmailRequest
	UpdateMeEmailResponse
*/
package api_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GetMeEmailRequest struct {
}

func (m *GetMeEmailRequest) Reset()                    { *m = GetMeEmailRequest{} }
func (m *GetMeEmailRequest) String() string            { return proto.CompactTextString(m) }
func (*GetMeEmailRequest) ProtoMessage()               {}
func (*GetMeEmailRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type GetMeEmailResponse struct {
	Email string `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
}

func (m *GetMeEmailResponse) Reset()                    { *m = GetMeEmailResponse{} }
func (m *GetMeEmailResponse) String() string            { return proto.CompactTextString(m) }
func (*GetMeEmailResponse) ProtoMessage()               {}
func (*GetMeEmailResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetMeEmailResponse) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type UpdateMeEmailRequest struct {
	Email string `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
}

func (m *UpdateMeEmailRequest) Reset()                    { *m = UpdateMeEmailRequest{} }
func (m *UpdateMeEmailRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateMeEmailRequest) ProtoMessage()               {}
func (*UpdateMeEmailRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *UpdateMeEmailRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type UpdateMeEmailResponse struct {
}

func (m *UpdateMeEmailResponse) Reset()                    { *m = UpdateMeEmailResponse{} }
func (m *UpdateMeEmailResponse) String() string            { return proto.CompactTextString(m) }
func (*UpdateMeEmailResponse) ProtoMessage()               {}
func (*UpdateMeEmailResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*GetMeEmailRequest)(nil), "api.v1.GetMeEmailRequest")
	proto.RegisterType((*GetMeEmailResponse)(nil), "api.v1.GetMeEmailResponse")
	proto.RegisterType((*UpdateMeEmailRequest)(nil), "api.v1.UpdateMeEmailRequest")
	proto.RegisterType((*UpdateMeEmailResponse)(nil), "api.v1.UpdateMeEmailResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for API service

type APIClient interface {
	GetMeEmail(ctx context.Context, in *GetMeEmailRequest, opts ...grpc.CallOption) (*GetMeEmailResponse, error)
	UpdateMeEmail(ctx context.Context, in *UpdateMeEmailRequest, opts ...grpc.CallOption) (*UpdateMeEmailResponse, error)
}

type aPIClient struct {
	cc *grpc.ClientConn
}

func NewAPIClient(cc *grpc.ClientConn) APIClient {
	return &aPIClient{cc}
}

func (c *aPIClient) GetMeEmail(ctx context.Context, in *GetMeEmailRequest, opts ...grpc.CallOption) (*GetMeEmailResponse, error) {
	out := new(GetMeEmailResponse)
	err := grpc.Invoke(ctx, "/api.v1.API/GetMeEmail", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) UpdateMeEmail(ctx context.Context, in *UpdateMeEmailRequest, opts ...grpc.CallOption) (*UpdateMeEmailResponse, error) {
	out := new(UpdateMeEmailResponse)
	err := grpc.Invoke(ctx, "/api.v1.API/UpdateMeEmail", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for API service

type APIServer interface {
	GetMeEmail(context.Context, *GetMeEmailRequest) (*GetMeEmailResponse, error)
	UpdateMeEmail(context.Context, *UpdateMeEmailRequest) (*UpdateMeEmailResponse, error)
}

func RegisterAPIServer(s *grpc.Server, srv APIServer) {
	s.RegisterService(&_API_serviceDesc, srv)
}

func _API_GetMeEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMeEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).GetMeEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.API/GetMeEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).GetMeEmail(ctx, req.(*GetMeEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_UpdateMeEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMeEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).UpdateMeEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.API/UpdateMeEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).UpdateMeEmail(ctx, req.(*UpdateMeEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _API_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.API",
	HandlerType: (*APIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMeEmail",
			Handler:    _API_GetMeEmail_Handler,
		},
		{
			MethodName: "UpdateMeEmail",
			Handler:    _API_UpdateMeEmail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/api/v1/me.proto",
}

func init() { proto.RegisterFile("proto/api/v1/me.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 225 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0x2c, 0xc8, 0xd4, 0x2f, 0x33, 0xd4, 0xcf, 0x4d, 0xd5, 0x03, 0xf3, 0x85, 0xd8,
	0x12, 0x0b, 0x32, 0xf5, 0xca, 0x0c, 0xa5, 0x64, 0xd2, 0xf3, 0xf3, 0xd3, 0x73, 0x52, 0xc1, 0xf2,
	0x89, 0x79, 0x79, 0xf9, 0x25, 0x89, 0x25, 0x99, 0xf9, 0x79, 0xc5, 0x10, 0x55, 0x4a, 0xc2, 0x5c,
	0x82, 0xee, 0xa9, 0x25, 0xbe, 0xa9, 0xae, 0xb9, 0x89, 0x99, 0x39, 0x41, 0xa9, 0x85, 0xa5, 0xa9,
	0xc5, 0x25, 0x4a, 0x5a, 0x5c, 0x42, 0xc8, 0x82, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x42, 0x22,
	0x5c, 0xac, 0xa9, 0x20, 0x01, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x08, 0x47, 0x49, 0x87,
	0x4b, 0x24, 0xb4, 0x20, 0x25, 0xb1, 0x24, 0x15, 0xd5, 0x0c, 0x1c, 0xaa, 0xc5, 0xb9, 0x44, 0xd1,
	0x54, 0x43, 0x0c, 0x37, 0x3a, 0xcf, 0xc8, 0xc5, 0xec, 0x18, 0xe0, 0x29, 0x14, 0xcb, 0xc5, 0x85,
	0xb0, 0x5a, 0x48, 0x52, 0x0f, 0xe2, 0x09, 0x3d, 0x0c, 0x37, 0x4a, 0x49, 0x61, 0x93, 0x82, 0x18,
	0xa6, 0x24, 0xd1, 0x74, 0xf9, 0xc9, 0x64, 0x26, 0x21, 0x21, 0x01, 0x44, 0xa0, 0xe8, 0x83, 0xed,
	0x17, 0xca, 0xe4, 0xe2, 0x45, 0xb1, 0x5f, 0x48, 0x06, 0x66, 0x0c, 0x36, 0x4f, 0x48, 0xc9, 0xe2,
	0x90, 0x85, 0xda, 0x23, 0x0d, 0xb6, 0x47, 0x54, 0x09, 0xc3, 0x1e, 0x2b, 0x46, 0xad, 0x24, 0x36,
	0x70, 0x00, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x1c, 0x87, 0x11, 0x47, 0x9f, 0x01, 0x00,
	0x00,
}
