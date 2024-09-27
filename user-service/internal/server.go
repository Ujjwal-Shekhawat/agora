package internal

import (
	"context"
	"log"
	"net"
	proto "proto/user"
	"user_service/config"
	"user_service/db"

	"github.com/gocql/gocql"
	"google.golang.org/grpc"
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

	var check_query string = "SELECT user_name, user_email from users where user_name = ?"
	user := map[string]interface{}{}
	if err := db.ExecQueryWithResponse(user, check_query, name); err != gocql.ErrNotFound && err != nil {
		response_proto.Message = "Error creating user"
		response_proto.StatusCode = 13
		log.Println("Exec error", err.Error())
		return response_proto, err
	}

	if len(user) > 0 {
		response_proto.Message = "UserName already exsists"
		response_proto.StatusCode = 13
		return response_proto, nil
	}

	var query string = "INSERT INTO users (user_id, user_name, user_email, user_password) VALUES (uuid(), ?, ?, ?);"
	if err := db.ExecQuery(query, name, email, password); err != nil {
		response_proto.Message = "Error creating user"
		response_proto.StatusCode = 13
		log.Fatal("Exec error", err.Error())
		return response_proto, err
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
