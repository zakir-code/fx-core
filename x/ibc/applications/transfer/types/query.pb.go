// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fx/ibc/applications/transfer/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	types "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func init() {
	proto.RegisterFile("fx/ibc/applications/transfer/v1/query.proto", fileDescriptor_569f08cc402420ba)
}

var fileDescriptor_569f08cc402420ba = []byte{
	// 305 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0x3f, 0x4b, 0xc4, 0x30,
	0x1c, 0x86, 0x5b, 0xd0, 0x03, 0x23, 0x2e, 0x19, 0x6f, 0x88, 0xe0, 0x24, 0x88, 0x89, 0xf5, 0xce,
	0x3f, 0x93, 0xa0, 0x28, 0x38, 0xaa, 0x38, 0x88, 0x5b, 0x9a, 0x4b, 0xbd, 0x82, 0x6d, 0x7a, 0xf9,
	0xa5, 0x67, 0x0f, 0xc1, 0xc9, 0x0f, 0xe0, 0xc7, 0x72, 0xbc, 0xd1, 0x51, 0xda, 0xef, 0xe0, 0x2c,
	0xad, 0xa4, 0xe7, 0x2d, 0x47, 0x33, 0x26, 0xbc, 0xcf, 0xfb, 0xbc, 0x43, 0x82, 0xf6, 0xa2, 0x82,
	0xc5, 0xa1, 0x60, 0x3c, 0xcb, 0x9e, 0x63, 0xc1, 0x4d, 0xac, 0x52, 0x60, 0x46, 0xf3, 0x14, 0x22,
	0xa9, 0xd9, 0x34, 0x60, 0x93, 0x5c, 0xea, 0x19, 0xcd, 0xb4, 0x32, 0x0a, 0x6f, 0x47, 0x05, 0x8d,
	0x43, 0x41, 0xff, 0x87, 0xa9, 0x0d, 0xd3, 0x69, 0xd0, 0xdf, 0xed, 0x5a, 0x75, 0xf8, 0xb3, 0x86,
	0xd6, 0x6f, 0xeb, 0x33, 0x7e, 0x45, 0xe8, 0x52, 0xa6, 0x2a, 0xb9, 0xd7, 0x5c, 0x48, 0x3c, 0x5c,
	0x29, 0xa0, 0x0d, 0xb2, 0x88, 0xdf, 0xc9, 0x49, 0x2e, 0xc1, 0xf4, 0x8f, 0x1c, 0x29, 0xc8, 0x54,
	0x0a, 0x72, 0xc7, 0xc3, 0x6f, 0x68, 0x73, 0x71, 0x0f, 0xd8, 0xad, 0x07, 0xac, 0xfe, 0xd8, 0x15,
	0x6b, 0xfd, 0x0a, 0xf5, 0x6e, 0xb8, 0xe6, 0x09, 0xe0, 0x83, 0x0e, 0x1d, 0x7f, 0x51, 0x6b, 0x0d,
	0x1c, 0x88, 0x56, 0x58, 0xa0, 0x8d, 0x66, 0xc9, 0x35, 0x87, 0x31, 0x1e, 0x74, 0xdd, 0x5d, 0xa7,
	0xad, 0x76, 0xe8, 0x06, 0xb5, 0xe6, 0x77, 0x1f, 0x6d, 0x5d, 0x81, 0xd0, 0xea, 0xe5, 0x7c, 0x34,
	0xd2, 0x12, 0x00, 0x9f, 0x74, 0x68, 0x5a, 0x22, 0xec, 0x84, 0x53, 0x77, 0xd0, 0xce, 0xb8, 0x78,
	0xf8, 0x2c, 0x89, 0x3f, 0x2f, 0x89, 0xff, 0x5d, 0x12, 0xff, 0xa3, 0x22, 0xde, 0xbc, 0x22, 0xde,
	0x57, 0x45, 0xbc, 0xc7, 0xb3, 0xa7, 0xd8, 0x8c, 0xf3, 0x90, 0x0a, 0x95, 0xb0, 0x28, 0x4f, 0x45,
	0x5d, 0x5c, 0xb0, 0xa8, 0xd8, 0x17, 0x4a, 0x4b, 0xb6, 0xea, 0x9b, 0x98, 0x59, 0x26, 0x21, 0xec,
	0x35, 0x2f, 0x7b, 0xf0, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xd2, 0xad, 0x42, 0x40, 0x53, 0x03, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// DenomTrace queries a denomination trace information.
	DenomTrace(ctx context.Context, in *types.QueryDenomTraceRequest, opts ...grpc.CallOption) (*types.QueryDenomTraceResponse, error)
	// DenomTraces queries all denomination traces.
	DenomTraces(ctx context.Context, in *types.QueryDenomTracesRequest, opts ...grpc.CallOption) (*types.QueryDenomTracesResponse, error)
	// Params queries all parameters of the ibc-transfer module.
	Params(ctx context.Context, in *types.QueryParamsRequest, opts ...grpc.CallOption) (*types.QueryParamsResponse, error)
	// DenomHash queries a denomination hash information.
	DenomHash(ctx context.Context, in *types.QueryDenomHashRequest, opts ...grpc.CallOption) (*types.QueryDenomHashResponse, error)
	// EscrowAddress returns the escrow address for a particular port and channel
	// id.
	EscrowAddress(ctx context.Context, in *types.QueryEscrowAddressRequest, opts ...grpc.CallOption) (*types.QueryEscrowAddressResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) DenomTrace(ctx context.Context, in *types.QueryDenomTraceRequest, opts ...grpc.CallOption) (*types.QueryDenomTraceResponse, error) {
	out := new(types.QueryDenomTraceResponse)
	err := c.cc.Invoke(ctx, "/fx.ibc.applications.transfer.v1.Query/DenomTrace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DenomTraces(ctx context.Context, in *types.QueryDenomTracesRequest, opts ...grpc.CallOption) (*types.QueryDenomTracesResponse, error) {
	out := new(types.QueryDenomTracesResponse)
	err := c.cc.Invoke(ctx, "/fx.ibc.applications.transfer.v1.Query/DenomTraces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Params(ctx context.Context, in *types.QueryParamsRequest, opts ...grpc.CallOption) (*types.QueryParamsResponse, error) {
	out := new(types.QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/fx.ibc.applications.transfer.v1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DenomHash(ctx context.Context, in *types.QueryDenomHashRequest, opts ...grpc.CallOption) (*types.QueryDenomHashResponse, error) {
	out := new(types.QueryDenomHashResponse)
	err := c.cc.Invoke(ctx, "/fx.ibc.applications.transfer.v1.Query/DenomHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) EscrowAddress(ctx context.Context, in *types.QueryEscrowAddressRequest, opts ...grpc.CallOption) (*types.QueryEscrowAddressResponse, error) {
	out := new(types.QueryEscrowAddressResponse)
	err := c.cc.Invoke(ctx, "/fx.ibc.applications.transfer.v1.Query/EscrowAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// DenomTrace queries a denomination trace information.
	DenomTrace(context.Context, *types.QueryDenomTraceRequest) (*types.QueryDenomTraceResponse, error)
	// DenomTraces queries all denomination traces.
	DenomTraces(context.Context, *types.QueryDenomTracesRequest) (*types.QueryDenomTracesResponse, error)
	// Params queries all parameters of the ibc-transfer module.
	Params(context.Context, *types.QueryParamsRequest) (*types.QueryParamsResponse, error)
	// DenomHash queries a denomination hash information.
	DenomHash(context.Context, *types.QueryDenomHashRequest) (*types.QueryDenomHashResponse, error)
	// EscrowAddress returns the escrow address for a particular port and channel
	// id.
	EscrowAddress(context.Context, *types.QueryEscrowAddressRequest) (*types.QueryEscrowAddressResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) DenomTrace(ctx context.Context, req *types.QueryDenomTraceRequest) (*types.QueryDenomTraceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DenomTrace not implemented")
}
func (*UnimplementedQueryServer) DenomTraces(ctx context.Context, req *types.QueryDenomTracesRequest) (*types.QueryDenomTracesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DenomTraces not implemented")
}
func (*UnimplementedQueryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) DenomHash(ctx context.Context, req *types.QueryDenomHashRequest) (*types.QueryDenomHashResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DenomHash not implemented")
}
func (*UnimplementedQueryServer) EscrowAddress(ctx context.Context, req *types.QueryEscrowAddressRequest) (*types.QueryEscrowAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EscrowAddress not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_DenomTrace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.QueryDenomTraceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DenomTrace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.ibc.applications.transfer.v1.Query/DenomTrace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DenomTrace(ctx, req.(*types.QueryDenomTraceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DenomTraces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.QueryDenomTracesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DenomTraces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.ibc.applications.transfer.v1.Query/DenomTraces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DenomTraces(ctx, req.(*types.QueryDenomTracesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.ibc.applications.transfer.v1.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*types.QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DenomHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.QueryDenomHashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DenomHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.ibc.applications.transfer.v1.Query/DenomHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DenomHash(ctx, req.(*types.QueryDenomHashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_EscrowAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.QueryEscrowAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).EscrowAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fx.ibc.applications.transfer.v1.Query/EscrowAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).EscrowAddress(ctx, req.(*types.QueryEscrowAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fx.ibc.applications.transfer.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DenomTrace",
			Handler:    _Query_DenomTrace_Handler,
		},
		{
			MethodName: "DenomTraces",
			Handler:    _Query_DenomTraces_Handler,
		},
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "DenomHash",
			Handler:    _Query_DenomHash_Handler,
		},
		{
			MethodName: "EscrowAddress",
			Handler:    _Query_EscrowAddress_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fx/ibc/applications/transfer/v1/query.proto",
}
