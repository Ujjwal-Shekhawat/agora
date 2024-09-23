package internal

import (
	"context"
	"log"
	"net"
	proto "proto/user"
	"user_service/config"
	"user_service/db"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, in *proto.GetUserReq) (*proto.CreateUserResponse, error) {
	return &proto.CreateUserResponse{
		Message:    "User id: " + in.UserID,
		StatusCode: 0,
	}, nil
}

func (s *server) CreateUser(ctx context.Context, in *proto.User) (*proto.CreateUserResponse, error) {
	response_proto := &proto.CreateUserResponse{
		Message:    "",
		StatusCode: 0,
	}

	name, email := in.Name, in.Email

	session, err := db.DatabaseSession()
	if err != nil {
		response_proto.Message = "Error creating user"
		response_proto.StatusCode = 13
		log.Println("Connection error", err.Error())
		return response_proto, err
	}
	defer session.Close()

	// Make a call to database to store user details (possibly implement db logic somewhere else please)
	var query string = "INSERT INTO users (user_id, user_name, user_email) VALUES (uuid(), ?, ?);"
	if err := session.Query(query, name, email).Exec(); err != nil {
		response_proto.Message = "Error creating user"
		response_proto.StatusCode = 13
		log.Println("Exec error", err.Error())
		return response_proto, err
	}

	response_proto.Message = "User created successfully"
	response_proto.StatusCode = 0

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
