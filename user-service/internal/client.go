package internal

import (
	"context"
	"log"
	proto "proto/user"
	"user_service/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type user_server struct {
	proto.UnimplementedUserServiceServer
}

func (s *user_server) GetUser(ctx context.Context, in *proto.GetUserReq) (*proto.UserResponse, error) {

	log.Println("Finding user", in.Name)
	res := db.MongoGetUser(in.Name)

	if len(res) == 0 {
		return nil, status.Errorf(codes.NotFound, "User not found with name %s", in.Name)
	}

	joinedGuilds := []string{}
	guilds, ok := res["guildNames"].(primitive.A)
	if ok {
		for _, v := range guilds {
			guild, ok := v.(string)
			if !ok {
				return nil, status.Error(codes.Internal, "Internal Server Error")
			}
			joinedGuilds = append(joinedGuilds, guild)
		}
	} else {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	return &proto.UserResponse{
		Name:         in.Name,
		JoinedGuilds: joinedGuilds,
	}, nil
}

func (s *user_server) CreateUser(ctx context.Context, in *proto.User) (*proto.ServerResponse, error) {
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

func (s *user_server) Login(ctx context.Context, loginReq *proto.LoginReq) (*proto.ServerResponse, error) {
	response_proto := &proto.ServerResponse{
		Message:    "",
		StatusCode: 0,
	}

	res := db.MongoGetUser(loginReq.Name)

	if len(res) == 0 {
		return nil, status.Error(codes.NotFound, "invalid credentials")
	} else {
		password := res["password"].(string)
		if password != loginReq.Password {
			return nil, status.Error(codes.PermissionDenied, "invalid credentials")
		}
	}

	response_proto.Message = "successfuly logged in"
	response_proto.StatusCode = 0

	return response_proto, nil
}
