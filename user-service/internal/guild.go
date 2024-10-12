package internal

import (
	"context"
	"log"
	proto "proto/guild"
	"time"
	"user_service/db"

	"github.com/gogo/protobuf/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	guildName, creator := guild.Name, guild.Creator

	if err := db.MongoCreateGuild(guildName, creator); err != nil {
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
		return nil, status.Error(codes.Internal, "Guild with that name does not exists")
	}

	if len(res) == 0 {
		log.Print("Guild not found here")
		return nil, status.Error(codes.NotFound, "Guild with that name does not exists")
	}

	channels := []string{}
	k, ok := res["channels"].(primitive.A)
	if ok {
		for _, channel := range k {
			name, ok := channel.(string)
			if ok {
				channels = append(channels, name)
			} else {
				log.Println(err)
				return nil, status.Error(codes.Internal, "Something went wrong terribly")
			}
		}
	} else {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Something went wrong terribly")
	}

	users := []string{}
	kk, ok := res["users"].(primitive.A)
	if ok {
		for _, channel := range kk {
			name, ok := channel.(string)
			if ok {
				users = append(users, name)
			} else {
				log.Println(err)
				return nil, status.Error(codes.Internal, "Something went wrong terribly")
			}
		}
	} else {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Something went wrong terribly")
	}

	response_proto.Name = res["name"].(string)
	response_proto.Members = users
	response_proto.Channels = channels

	return response_proto, nil
}

func (s *guild_server) JoinGuild(ctx context.Context, guild *proto.GuildMember) (*proto.ServerResponse, error) {
	response_proto := &proto.ServerResponse{
		Message:    "",
		StatusCode: 0,
	}

	memberName, guildName := guild.Name, guild.GuildName

	log.Println(memberName, guildName)

	if err := db.MongoJoinGuild(memberName, guildName); err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Something went wrong while adding user guild")
	}

	response_proto.Message = "Joined guild successfully"
	response_proto.StatusCode = 0

	return response_proto, nil
}

func (s *guild_server) LeaveGuild(ctx context.Context, guild *proto.GuildMember) (*proto.ServerResponse, error) {
	return nil, nil
}

func (s *guild_server) GetMessages(ctx context.Context, messageRequest *proto.GuildMessagesRequest) (*proto.GuildMessagesResponse, error) {
	response_proto := &proto.GuildMessagesResponse{
		Messages: []*proto.GuildMessageFormat{},
	}

	messages := &[]struct {
		UserName  string
		Message   string
		Timestamp time.Time
	}{}
	err := db.ExecQueryWithResponse("select user_name, user_message, timestamp from messages where guild_name = 'guild-FairyTale' and channel_name='general' order by timestamp desc limit 100;", messages)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Aborted, "Well some error occured")
	}

	for _, message := range *messages {
		t, err := types.TimestampProto(message.Timestamp)
		if err != nil {
			return nil, status.Error(codes.Internal, "Timestamp issue please check")
		}
		response_proto.Messages = append(response_proto.Messages, &proto.GuildMessageFormat{
			Key:       message.UserName,
			Value:     message.Message,
			Timestamp: t,
		})
	}

	return response_proto, nil
}
