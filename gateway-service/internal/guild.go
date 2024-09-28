package internal

import (
	"context"
	gproto "proto/guild"
	"time"
)

func (c *ServiceClientStruct) MakeGuild(guild *gproto.Guild) (*gproto.ServerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &gproto.Guild{Name: guild.Name}
	res, err := c.guild.CreateGuild(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
