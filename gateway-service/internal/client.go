package internal

import (
	"context"
	"log"
	proto "proto/user"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClientStruct struct {
	client proto.UserServiceClient
}

func GetUserServiceClient(addr string) (*UserServiceClientStruct, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &UserServiceClientStruct{
		client: proto.NewUserServiceClient(conn),
	}, nil
}

func (c *UserServiceClientStruct) GetUserDetails(name string) (*proto.ServerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &proto.GetUserReq{Name: name}
	res, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *UserServiceClientStruct) CreateNewUser(user *proto.User) (*proto.ServerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println(user.Name, user.Email, user.Password)

	req := &proto.User{Name: user.Name, Email: user.Email, Password: user.Password}
	res, err := c.client.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *UserServiceClientStruct) LoginUser(login *proto.LoginReq) (*proto.ServerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &proto.LoginReq{Name: login.Name, Password: login.Password}
	res, err := c.client.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
