package internal

import (
	"context"
	proto "proto/user"
	"time"
	"user_service/models"

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

func (c *UserServiceClientStruct) GetUserDetails(userId string) (*proto.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &proto.GetUserReq{UserID: userId}
	res, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *UserServiceClientStruct) CreateNewUser(user *models.User) (*proto.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &proto.User{Name: user.Name, Email: user.Email}
	res, err := c.client.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
