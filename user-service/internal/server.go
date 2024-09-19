package internal

import (
	"context"
	"net"
	proto "proto/user"
	"user_service/config"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, in *proto.GetUserReq) (*proto.CreateUserResponse, error) {
	return &proto.CreateUserResponse{
		Message:    "User id: " + in.UserID,
		StatusCode: 200,
	}, nil
}

func StartGrpcServer(cfg *config.ServiceConfig) error {
	lis, err := net.Listen("tcp", cfg.ServerPort)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
