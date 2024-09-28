package internal

import (
	"context"
	"log"
	proto "proto/guild"
	"user_service/db"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type guild_server struct {
	proto.UnimplementedGuildServiceServer
}

func (s *guild_server) CreateGuild(ctx context.Context, guild *proto.Guild) (*proto.ServerResponse, error) {
	response_proto := &proto.ServerResponse{
		Message:    "",
		StatusCode: 0,
	}

	guildName := guild.Name

	if err := db.MongoCreateGuild(guildName); err != nil {
		log.Println(err)
		return nil, status.Error(codes.InvalidArgument, "Guild with that name already exsists")
	}

	response_proto.Message = "Guild created successfully"
	response_proto.StatusCode = 0

	return response_proto, nil
}

func (s *guild_server) GetGuild(ctx context.Context, guild *proto.Guild) (*proto.GuildResponse, error) {
	response_proto := &proto.GuildResponse{
		Name:     "",
		Channels: []string{""},
	}

	guildName := guild.Name

	res, err := db.MongoGetGuild(guildName)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Something went wrong terribly")
	}

	if len(res) == 0 {
		log.Print("Guild not found here")
		return nil, status.Error(codes.NotFound, "Guild with that name doesnot exists")
	}

	response_proto.Name = res["name"].(string)
	response_proto.Channels = res["channels"].([]string)

	return response_proto, nil
}

func (s *guild_server) JoinGuild(ctx context.Context, guild *proto.GuildMember) (*proto.ServerResponse, error) {
	return nil, nil
}

func (s *guild_server) LeaveGuild(ctx context.Context, guild *proto.GuildMember) (*proto.ServerResponse, error) {
	return nil, nil
}
