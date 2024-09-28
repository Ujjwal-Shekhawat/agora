package internal

import (
	"net"
	gproto "proto/guild"
	uproto "proto/user"
	"user_service/config"

	"google.golang.org/grpc"
)

func StartGrpcServer(cfg *config.ServiceConfig) error {
	lis, err := net.Listen("tcp", cfg.ServerPort)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	uproto.RegisterUserServiceServer(s, &user_server{})
	gproto.RegisterGuildServiceServer(s, &guild_server{})
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
