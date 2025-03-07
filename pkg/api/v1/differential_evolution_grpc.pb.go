// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: api/v1/differential_evolution.proto

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
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DifferentialEvolutionService_ListSupportedAlgorithms_FullMethodName = "/api.v1.DifferentialEvolutionService/ListSupportedAlgorithms"
	DifferentialEvolutionService_ListSupportedVariants_FullMethodName   = "/api.v1.DifferentialEvolutionService/ListSupportedVariants"
	DifferentialEvolutionService_ListSupportedProblems_FullMethodName   = "/api.v1.DifferentialEvolutionService/ListSupportedProblems"
	DifferentialEvolutionService_Run_FullMethodName                     = "/api.v1.DifferentialEvolutionService/Run"
)

// DifferentialEvolutionServiceClient is the client API for DifferentialEvolutionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DifferentialEvolutionServiceClient interface {
	ListSupportedAlgorithms(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListSupportedAlgorithmsResponse, error)
	ListSupportedVariants(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListSupportedVariantsResponse, error)
	ListSupportedProblems(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListSupportedProblemsResponse, error)
	Run(ctx context.Context, in *RunRequest, opts ...grpc.CallOption) (*RunResponse, error)
}

type differentialEvolutionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDifferentialEvolutionServiceClient(cc grpc.ClientConnInterface) DifferentialEvolutionServiceClient {
	return &differentialEvolutionServiceClient{cc}
}

func (c *differentialEvolutionServiceClient) ListSupportedAlgorithms(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListSupportedAlgorithmsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListSupportedAlgorithmsResponse)
	err := c.cc.Invoke(ctx, DifferentialEvolutionService_ListSupportedAlgorithms_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *differentialEvolutionServiceClient) ListSupportedVariants(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListSupportedVariantsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListSupportedVariantsResponse)
	err := c.cc.Invoke(ctx, DifferentialEvolutionService_ListSupportedVariants_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *differentialEvolutionServiceClient) ListSupportedProblems(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListSupportedProblemsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListSupportedProblemsResponse)
	err := c.cc.Invoke(ctx, DifferentialEvolutionService_ListSupportedProblems_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *differentialEvolutionServiceClient) Run(ctx context.Context, in *RunRequest, opts ...grpc.CallOption) (*RunResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RunResponse)
	err := c.cc.Invoke(ctx, DifferentialEvolutionService_Run_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DifferentialEvolutionServiceServer is the server API for DifferentialEvolutionService service.
// All implementations must embed UnimplementedDifferentialEvolutionServiceServer
// for forward compatibility.
type DifferentialEvolutionServiceServer interface {
	ListSupportedAlgorithms(context.Context, *emptypb.Empty) (*ListSupportedAlgorithmsResponse, error)
	ListSupportedVariants(context.Context, *emptypb.Empty) (*ListSupportedVariantsResponse, error)
	ListSupportedProblems(context.Context, *emptypb.Empty) (*ListSupportedProblemsResponse, error)
	Run(context.Context, *RunRequest) (*RunResponse, error)
	mustEmbedUnimplementedDifferentialEvolutionServiceServer()
}

// UnimplementedDifferentialEvolutionServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDifferentialEvolutionServiceServer struct{}

func (UnimplementedDifferentialEvolutionServiceServer) ListSupportedAlgorithms(context.Context, *emptypb.Empty) (*ListSupportedAlgorithmsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSupportedAlgorithms not implemented")
}
func (UnimplementedDifferentialEvolutionServiceServer) ListSupportedVariants(context.Context, *emptypb.Empty) (*ListSupportedVariantsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSupportedVariants not implemented")
}
func (UnimplementedDifferentialEvolutionServiceServer) ListSupportedProblems(context.Context, *emptypb.Empty) (*ListSupportedProblemsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSupportedProblems not implemented")
}
func (UnimplementedDifferentialEvolutionServiceServer) Run(context.Context, *RunRequest) (*RunResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Run not implemented")
}
func (UnimplementedDifferentialEvolutionServiceServer) mustEmbedUnimplementedDifferentialEvolutionServiceServer() {
}
func (UnimplementedDifferentialEvolutionServiceServer) testEmbeddedByValue() {}

// UnsafeDifferentialEvolutionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DifferentialEvolutionServiceServer will
// result in compilation errors.
type UnsafeDifferentialEvolutionServiceServer interface {
	mustEmbedUnimplementedDifferentialEvolutionServiceServer()
}

func RegisterDifferentialEvolutionServiceServer(s grpc.ServiceRegistrar, srv DifferentialEvolutionServiceServer) {
	// If the following call pancis, it indicates UnimplementedDifferentialEvolutionServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DifferentialEvolutionService_ServiceDesc, srv)
}

func _DifferentialEvolutionService_ListSupportedAlgorithms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DifferentialEvolutionServiceServer).ListSupportedAlgorithms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DifferentialEvolutionService_ListSupportedAlgorithms_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DifferentialEvolutionServiceServer).ListSupportedAlgorithms(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DifferentialEvolutionService_ListSupportedVariants_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DifferentialEvolutionServiceServer).ListSupportedVariants(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DifferentialEvolutionService_ListSupportedVariants_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DifferentialEvolutionServiceServer).ListSupportedVariants(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DifferentialEvolutionService_ListSupportedProblems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DifferentialEvolutionServiceServer).ListSupportedProblems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DifferentialEvolutionService_ListSupportedProblems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DifferentialEvolutionServiceServer).ListSupportedProblems(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DifferentialEvolutionService_Run_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DifferentialEvolutionServiceServer).Run(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DifferentialEvolutionService_Run_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DifferentialEvolutionServiceServer).Run(ctx, req.(*RunRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DifferentialEvolutionService_ServiceDesc is the grpc.ServiceDesc for DifferentialEvolutionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DifferentialEvolutionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.DifferentialEvolutionService",
	HandlerType: (*DifferentialEvolutionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListSupportedAlgorithms",
			Handler:    _DifferentialEvolutionService_ListSupportedAlgorithms_Handler,
		},
		{
			MethodName: "ListSupportedVariants",
			Handler:    _DifferentialEvolutionService_ListSupportedVariants_Handler,
		},
		{
			MethodName: "ListSupportedProblems",
			Handler:    _DifferentialEvolutionService_ListSupportedProblems_Handler,
		},
		{
			MethodName: "Run",
			Handler:    _DifferentialEvolutionService_Run_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/differential_evolution.proto",
}
