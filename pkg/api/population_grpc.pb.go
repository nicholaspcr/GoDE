// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: api/population.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PopulationServices_Create_FullMethodName     = "/api.PopulationServices/Create"
	PopulationServices_Read_FullMethodName       = "/api.PopulationServices/Read"
	PopulationServices_Update_FullMethodName     = "/api.PopulationServices/Update"
	PopulationServices_Delete_FullMethodName     = "/api.PopulationServices/Delete"
	PopulationServices_ListByUser_FullMethodName = "/api.PopulationServices/ListByUser"
)

// PopulationServicesClient is the client API for PopulationServices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PopulationServicesClient interface {
	Create(ctx context.Context, in *Population, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Read(ctx context.Context, in *PopulationIDs, opts ...grpc.CallOption) (*Population, error)
	Update(ctx context.Context, in *Population, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Delete(ctx context.Context, in *PopulationIDs, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ListByUser(ctx context.Context, in *UserIDs, opts ...grpc.CallOption) (PopulationServices_ListByUserClient, error)
}

type populationServicesClient struct {
	cc grpc.ClientConnInterface
}

func NewPopulationServicesClient(cc grpc.ClientConnInterface) PopulationServicesClient {
	return &populationServicesClient{cc}
}

func (c *populationServicesClient) Create(ctx context.Context, in *Population, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PopulationServices_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *populationServicesClient) Read(ctx context.Context, in *PopulationIDs, opts ...grpc.CallOption) (*Population, error) {
	out := new(Population)
	err := c.cc.Invoke(ctx, PopulationServices_Read_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *populationServicesClient) Update(ctx context.Context, in *Population, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PopulationServices_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *populationServicesClient) Delete(ctx context.Context, in *PopulationIDs, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PopulationServices_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *populationServicesClient) ListByUser(ctx context.Context, in *UserIDs, opts ...grpc.CallOption) (PopulationServices_ListByUserClient, error) {
	stream, err := c.cc.NewStream(ctx, &PopulationServices_ServiceDesc.Streams[0], PopulationServices_ListByUser_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &populationServicesListByUserClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PopulationServices_ListByUserClient interface {
	Recv() (*Population, error)
	grpc.ClientStream
}

type populationServicesListByUserClient struct {
	grpc.ClientStream
}

func (x *populationServicesListByUserClient) Recv() (*Population, error) {
	m := new(Population)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PopulationServicesServer is the server API for PopulationServices service.
// All implementations must embed UnimplementedPopulationServicesServer
// for forward compatibility
type PopulationServicesServer interface {
	Create(context.Context, *Population) (*emptypb.Empty, error)
	Read(context.Context, *PopulationIDs) (*Population, error)
	Update(context.Context, *Population) (*emptypb.Empty, error)
	Delete(context.Context, *PopulationIDs) (*emptypb.Empty, error)
	ListByUser(*UserIDs, PopulationServices_ListByUserServer) error
	mustEmbedUnimplementedPopulationServicesServer()
}

// UnimplementedPopulationServicesServer must be embedded to have forward compatible implementations.
type UnimplementedPopulationServicesServer struct {
}

func (UnimplementedPopulationServicesServer) Create(context.Context, *Population) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedPopulationServicesServer) Read(context.Context, *PopulationIDs) (*Population, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedPopulationServicesServer) Update(context.Context, *Population) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedPopulationServicesServer) Delete(context.Context, *PopulationIDs) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedPopulationServicesServer) ListByUser(*UserIDs, PopulationServices_ListByUserServer) error {
	return status.Errorf(codes.Unimplemented, "method ListByUser not implemented")
}
func (UnimplementedPopulationServicesServer) mustEmbedUnimplementedPopulationServicesServer() {}

// UnsafePopulationServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PopulationServicesServer will
// result in compilation errors.
type UnsafePopulationServicesServer interface {
	mustEmbedUnimplementedPopulationServicesServer()
}

func RegisterPopulationServicesServer(s grpc.ServiceRegistrar, srv PopulationServicesServer) {
	s.RegisterService(&PopulationServices_ServiceDesc, srv)
}

func _PopulationServices_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Population)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopulationServicesServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PopulationServices_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopulationServicesServer).Create(ctx, req.(*Population))
	}
	return interceptor(ctx, in, info, handler)
}

func _PopulationServices_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PopulationIDs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopulationServicesServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PopulationServices_Read_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopulationServicesServer).Read(ctx, req.(*PopulationIDs))
	}
	return interceptor(ctx, in, info, handler)
}

func _PopulationServices_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Population)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopulationServicesServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PopulationServices_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopulationServicesServer).Update(ctx, req.(*Population))
	}
	return interceptor(ctx, in, info, handler)
}

func _PopulationServices_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PopulationIDs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopulationServicesServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PopulationServices_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopulationServicesServer).Delete(ctx, req.(*PopulationIDs))
	}
	return interceptor(ctx, in, info, handler)
}

func _PopulationServices_ListByUser_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(UserIDs)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PopulationServicesServer).ListByUser(m, &populationServicesListByUserServer{stream})
}

type PopulationServices_ListByUserServer interface {
	Send(*Population) error
	grpc.ServerStream
}

type populationServicesListByUserServer struct {
	grpc.ServerStream
}

func (x *populationServicesListByUserServer) Send(m *Population) error {
	return x.ServerStream.SendMsg(m)
}

// PopulationServices_ServiceDesc is the grpc.ServiceDesc for PopulationServices service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PopulationServices_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.PopulationServices",
	HandlerType: (*PopulationServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _PopulationServices_Create_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _PopulationServices_Read_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _PopulationServices_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _PopulationServices_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListByUser",
			Handler:       _PopulationServices_ListByUser_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/population.proto",
}

const (
	ParetoServices_Create_FullMethodName     = "/api.ParetoServices/Create"
	ParetoServices_Read_FullMethodName       = "/api.ParetoServices/Read"
	ParetoServices_Update_FullMethodName     = "/api.ParetoServices/Update"
	ParetoServices_Delete_FullMethodName     = "/api.ParetoServices/Delete"
	ParetoServices_ListByUser_FullMethodName = "/api.ParetoServices/ListByUser"
)

// ParetoServicesClient is the client API for ParetoServices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ParetoServicesClient interface {
	Create(ctx context.Context, in *Pareto, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Read(ctx context.Context, in *ParetoIDs, opts ...grpc.CallOption) (*Pareto, error)
	Update(ctx context.Context, in *Pareto, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Delete(ctx context.Context, in *ParetoIDs, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ListByUser(ctx context.Context, in *UserIDs, opts ...grpc.CallOption) (ParetoServices_ListByUserClient, error)
}

type paretoServicesClient struct {
	cc grpc.ClientConnInterface
}

func NewParetoServicesClient(cc grpc.ClientConnInterface) ParetoServicesClient {
	return &paretoServicesClient{cc}
}

func (c *paretoServicesClient) Create(ctx context.Context, in *Pareto, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ParetoServices_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paretoServicesClient) Read(ctx context.Context, in *ParetoIDs, opts ...grpc.CallOption) (*Pareto, error) {
	out := new(Pareto)
	err := c.cc.Invoke(ctx, ParetoServices_Read_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paretoServicesClient) Update(ctx context.Context, in *Pareto, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ParetoServices_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paretoServicesClient) Delete(ctx context.Context, in *ParetoIDs, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ParetoServices_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paretoServicesClient) ListByUser(ctx context.Context, in *UserIDs, opts ...grpc.CallOption) (ParetoServices_ListByUserClient, error) {
	stream, err := c.cc.NewStream(ctx, &ParetoServices_ServiceDesc.Streams[0], ParetoServices_ListByUser_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &paretoServicesListByUserClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ParetoServices_ListByUserClient interface {
	Recv() (*Pareto, error)
	grpc.ClientStream
}

type paretoServicesListByUserClient struct {
	grpc.ClientStream
}

func (x *paretoServicesListByUserClient) Recv() (*Pareto, error) {
	m := new(Pareto)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ParetoServicesServer is the server API for ParetoServices service.
// All implementations must embed UnimplementedParetoServicesServer
// for forward compatibility
type ParetoServicesServer interface {
	Create(context.Context, *Pareto) (*emptypb.Empty, error)
	Read(context.Context, *ParetoIDs) (*Pareto, error)
	Update(context.Context, *Pareto) (*emptypb.Empty, error)
	Delete(context.Context, *ParetoIDs) (*emptypb.Empty, error)
	ListByUser(*UserIDs, ParetoServices_ListByUserServer) error
	mustEmbedUnimplementedParetoServicesServer()
}

// UnimplementedParetoServicesServer must be embedded to have forward compatible implementations.
type UnimplementedParetoServicesServer struct {
}

func (UnimplementedParetoServicesServer) Create(context.Context, *Pareto) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedParetoServicesServer) Read(context.Context, *ParetoIDs) (*Pareto, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedParetoServicesServer) Update(context.Context, *Pareto) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedParetoServicesServer) Delete(context.Context, *ParetoIDs) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedParetoServicesServer) ListByUser(*UserIDs, ParetoServices_ListByUserServer) error {
	return status.Errorf(codes.Unimplemented, "method ListByUser not implemented")
}
func (UnimplementedParetoServicesServer) mustEmbedUnimplementedParetoServicesServer() {}

// UnsafeParetoServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ParetoServicesServer will
// result in compilation errors.
type UnsafeParetoServicesServer interface {
	mustEmbedUnimplementedParetoServicesServer()
}

func RegisterParetoServicesServer(s grpc.ServiceRegistrar, srv ParetoServicesServer) {
	s.RegisterService(&ParetoServices_ServiceDesc, srv)
}

func _ParetoServices_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Pareto)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ParetoServicesServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ParetoServices_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ParetoServicesServer).Create(ctx, req.(*Pareto))
	}
	return interceptor(ctx, in, info, handler)
}

func _ParetoServices_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParetoIDs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ParetoServicesServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ParetoServices_Read_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ParetoServicesServer).Read(ctx, req.(*ParetoIDs))
	}
	return interceptor(ctx, in, info, handler)
}

func _ParetoServices_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Pareto)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ParetoServicesServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ParetoServices_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ParetoServicesServer).Update(ctx, req.(*Pareto))
	}
	return interceptor(ctx, in, info, handler)
}

func _ParetoServices_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParetoIDs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ParetoServicesServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ParetoServices_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ParetoServicesServer).Delete(ctx, req.(*ParetoIDs))
	}
	return interceptor(ctx, in, info, handler)
}

func _ParetoServices_ListByUser_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(UserIDs)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ParetoServicesServer).ListByUser(m, &paretoServicesListByUserServer{stream})
}

type ParetoServices_ListByUserServer interface {
	Send(*Pareto) error
	grpc.ServerStream
}

type paretoServicesListByUserServer struct {
	grpc.ServerStream
}

func (x *paretoServicesListByUserServer) Send(m *Pareto) error {
	return x.ServerStream.SendMsg(m)
}

// ParetoServices_ServiceDesc is the grpc.ServiceDesc for ParetoServices service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ParetoServices_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.ParetoServices",
	HandlerType: (*ParetoServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _ParetoServices_Create_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _ParetoServices_Read_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _ParetoServices_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ParetoServices_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListByUser",
			Handler:       _ParetoServices_ListByUser_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/population.proto",
}
