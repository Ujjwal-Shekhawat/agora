package internal

import (
	"context"
	"log"
	uproto "proto/user"
	"time"
)

func (c *ServiceClientStruct) GetUserDetails(name string) (*uproto.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &uproto.GetUserReq{Name: name}
	res, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ServiceClientStruct) CreateNewUser(user *uproto.User) (*uproto.ServerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println(user.Name, user.Email, user.Password)

	req := &uproto.User{Name: user.Name, Email: user.Email, Password: user.Password}
	res, err := c.client.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ServiceClientStruct) LoginUser(login *uproto.LoginReq) (*uproto.ServerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &uproto.LoginReq{Name: login.Name, Password: login.Password}
	res, err := c.client.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
