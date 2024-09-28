package internal

import (
	"context"
	"log"
	"net"
	proto "proto/user"
	"user_service/config"
	"user_service/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, in *proto.GetUserReq) (*proto.ServerResponse, error) {
	return &proto.ServerResponse{
		Message:    "User id: " + in.Name,
		StatusCode: 0,
	}, nil
}

func (s *server) CreateUser(ctx context.Context, in *proto.User) (*proto.ServerResponse, error) {
	response_proto := &proto.ServerResponse{
		Message:    "",
		StatusCode: 0,
	}

	name, email, password := in.Name, in.Email, in.Password

	err := db.MongoCreateUser(name, email, password)
	if err != nil {
		log.Println("User already exsists", err)
		return nil, status.Errorf(codes.AlreadyExists, "User already exsists")
	}

	response_proto.Message = "User created successfully"
	response_proto.StatusCode = 0

	return response_proto, nil
}

func (s *server) Login(ctx context.Context, loginReq *proto.LoginReq) (*proto.ServerResponse, error) {
	response_proto := &proto.ServerResponse{
		Message:    "",
		StatusCode: 0,
	}

	response_proto.Message = "Login not yet implemented"
	response_proto.StatusCode = 13

	return response_proto, nil
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
