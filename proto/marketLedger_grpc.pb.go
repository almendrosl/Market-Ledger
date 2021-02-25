// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package marketLedgerGrpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MarketLedgerServiceClient is the client API for MarketLedgerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MarketLedgerServiceClient interface {
	CreateInvoice(ctx context.Context, in *CreateInvoiceReq, opts ...grpc.CallOption) (*CreateInvoiceResp, error)
	SellOrders(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*SellOrdersResp, error)
}

type marketLedgerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMarketLedgerServiceClient(cc grpc.ClientConnInterface) MarketLedgerServiceClient {
	return &marketLedgerServiceClient{cc}
}

func (c *marketLedgerServiceClient) CreateInvoice(ctx context.Context, in *CreateInvoiceReq, opts ...grpc.CallOption) (*CreateInvoiceResp, error) {
	out := new(CreateInvoiceResp)
	err := c.cc.Invoke(ctx, "/marketLedgerGrpc.MarketLedgerService/CreateInvoice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketLedgerServiceClient) SellOrders(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*SellOrdersResp, error) {
	out := new(SellOrdersResp)
	err := c.cc.Invoke(ctx, "/marketLedgerGrpc.MarketLedgerService/SellOrders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MarketLedgerServiceServer is the server API for MarketLedgerService service.
// All implementations must embed UnimplementedMarketLedgerServiceServer
// for forward compatibility
type MarketLedgerServiceServer interface {
	CreateInvoice(context.Context, *CreateInvoiceReq) (*CreateInvoiceResp, error)
	SellOrders(context.Context, *Empty) (*SellOrdersResp, error)
	mustEmbedUnimplementedMarketLedgerServiceServer()
}

// UnimplementedMarketLedgerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMarketLedgerServiceServer struct {
}

func (UnimplementedMarketLedgerServiceServer) CreateInvoice(context.Context, *CreateInvoiceReq) (*CreateInvoiceResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateInvoice not implemented")
}
func (UnimplementedMarketLedgerServiceServer) SellOrders(context.Context, *Empty) (*SellOrdersResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SellOrders not implemented")
}
func (UnimplementedMarketLedgerServiceServer) mustEmbedUnimplementedMarketLedgerServiceServer() {}

// UnsafeMarketLedgerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MarketLedgerServiceServer will
// result in compilation errors.
type UnsafeMarketLedgerServiceServer interface {
	mustEmbedUnimplementedMarketLedgerServiceServer()
}

func RegisterMarketLedgerServiceServer(s grpc.ServiceRegistrar, srv MarketLedgerServiceServer) {
	s.RegisterService(&MarketLedgerService_ServiceDesc, srv)
}

func _MarketLedgerService_CreateInvoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateInvoiceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketLedgerServiceServer).CreateInvoice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/marketLedgerGrpc.MarketLedgerService/CreateInvoice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketLedgerServiceServer).CreateInvoice(ctx, req.(*CreateInvoiceReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MarketLedgerService_SellOrders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketLedgerServiceServer).SellOrders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/marketLedgerGrpc.MarketLedgerService/SellOrders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketLedgerServiceServer).SellOrders(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// MarketLedgerService_ServiceDesc is the grpc.ServiceDesc for MarketLedgerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MarketLedgerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "marketLedgerGrpc.MarketLedgerService",
	HandlerType: (*MarketLedgerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateInvoice",
			Handler:    _MarketLedgerService_CreateInvoice_Handler,
		},
		{
			MethodName: "SellOrders",
			Handler:    _MarketLedgerService_SellOrders_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/marketLedger.proto",
}
