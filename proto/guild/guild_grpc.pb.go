// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.28.2
// source: guild/guild.proto

package proto

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

// GuildServiceClient is the client API for GuildService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GuildServiceClient interface {
	CreateGuild(ctx context.Context, in *Guild, opts ...grpc.CallOption) (*ServerResponse, error)
	GetGuild(ctx context.Context, in *Guild, opts ...grpc.CallOption) (*GuildResponse, error)
	JoinGuild(ctx context.Context, in *GuildMember, opts ...grpc.CallOption) (*ServerResponse, error)
	LeaveGuild(ctx context.Context, in *GuildMember, opts ...grpc.CallOption) (*ServerResponse, error)
}

type guildServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGuildServiceClient(cc grpc.ClientConnInterface) GuildServiceClient {
	return &guildServiceClient{cc}
}

func (c *guildServiceClient) CreateGuild(ctx context.Context, in *Guild, opts ...grpc.CallOption) (*ServerResponse, error) {
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, "/guild_proto.GuildService/CreateGuild", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *guildServiceClient) GetGuild(ctx context.Context, in *Guild, opts ...grpc.CallOption) (*GuildResponse, error) {
	out := new(GuildResponse)
	err := c.cc.Invoke(ctx, "/guild_proto.GuildService/GetGuild", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *guildServiceClient) JoinGuild(ctx context.Context, in *GuildMember, opts ...grpc.CallOption) (*ServerResponse, error) {
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, "/guild_proto.GuildService/JoinGuild", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *guildServiceClient) LeaveGuild(ctx context.Context, in *GuildMember, opts ...grpc.CallOption) (*ServerResponse, error) {
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, "/guild_proto.GuildService/LeaveGuild", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GuildServiceServer is the server API for GuildService service.
// All implementations must embed UnimplementedGuildServiceServer
// for forward compatibility
type GuildServiceServer interface {
	CreateGuild(context.Context, *Guild) (*ServerResponse, error)
	GetGuild(context.Context, *Guild) (*GuildResponse, error)
	JoinGuild(context.Context, *GuildMember) (*ServerResponse, error)
	LeaveGuild(context.Context, *GuildMember) (*ServerResponse, error)
	mustEmbedUnimplementedGuildServiceServer()
}

// UnimplementedGuildServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGuildServiceServer struct {
}

func (UnimplementedGuildServiceServer) CreateGuild(context.Context, *Guild) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGuild not implemented")
}
func (UnimplementedGuildServiceServer) GetGuild(context.Context, *Guild) (*GuildResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGuild not implemented")
}
func (UnimplementedGuildServiceServer) JoinGuild(context.Context, *GuildMember) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinGuild not implemented")
}
func (UnimplementedGuildServiceServer) LeaveGuild(context.Context, *GuildMember) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveGuild not implemented")
}
func (UnimplementedGuildServiceServer) mustEmbedUnimplementedGuildServiceServer() {}

// UnsafeGuildServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GuildServiceServer will
// result in compilation errors.
type UnsafeGuildServiceServer interface {
	mustEmbedUnimplementedGuildServiceServer()
}

func RegisterGuildServiceServer(s grpc.ServiceRegistrar, srv GuildServiceServer) {
	s.RegisterService(&GuildService_ServiceDesc, srv)
}

func _GuildService_CreateGuild_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Guild)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GuildServiceServer).CreateGuild(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/guild_proto.GuildService/CreateGuild",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GuildServiceServer).CreateGuild(ctx, req.(*Guild))
	}
	return interceptor(ctx, in, info, handler)
}

func _GuildService_GetGuild_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Guild)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GuildServiceServer).GetGuild(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/guild_proto.GuildService/GetGuild",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GuildServiceServer).GetGuild(ctx, req.(*Guild))
	}
	return interceptor(ctx, in, info, handler)
}

func _GuildService_JoinGuild_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GuildMember)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GuildServiceServer).JoinGuild(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/guild_proto.GuildService/JoinGuild",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GuildServiceServer).JoinGuild(ctx, req.(*GuildMember))
	}
	return interceptor(ctx, in, info, handler)
}

func _GuildService_LeaveGuild_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GuildMember)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GuildServiceServer).LeaveGuild(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/guild_proto.GuildService/LeaveGuild",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GuildServiceServer).LeaveGuild(ctx, req.(*GuildMember))
	}
	return interceptor(ctx, in, info, handler)
}

// GuildService_ServiceDesc is the grpc.ServiceDesc for GuildService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GuildService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "guild_proto.GuildService",
	HandlerType: (*GuildServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateGuild",
			Handler:    _GuildService_CreateGuild_Handler,
		},
		{
			MethodName: "GetGuild",
			Handler:    _GuildService_GetGuild_Handler,
		},
		{
			MethodName: "JoinGuild",
			Handler:    _GuildService_JoinGuild_Handler,
		},
		{
			MethodName: "LeaveGuild",
			Handler:    _GuildService_LeaveGuild_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "guild/guild.proto",
}