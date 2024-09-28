package internal

import (
	gproto "proto/guild"
	uproto "proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClientStruct struct {
	client uproto.UserServiceClient
	guild  gproto.GuildServiceClient
}

func GetServiceClient(addr string) (*ServiceClientStruct, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ServiceClientStruct{
		client: uproto.NewUserServiceClient(conn),
		guild:  gproto.NewGuildServiceClient(conn),
	}, nil
}
