package actors

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/okzmo/kyob/db"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
	protoTypes "github.com/okzmo/kyob/types"
)

// SERVER

func (s *server) RemoveServer(ctx *actor.Context, msg *protoTypes.BodyServerRemoved) {
	err := services.DeleteServer(context.TODO(), msg.ServerId, msg.UserId)
	if err != nil {
		slog.Error("failed to delete server", "err", err)
		return
	}

	ctx.Engine().Poison(ctx.PID())
}

func (s *server) NewUser(ctx *actor.Context, msg *protoTypes.BodyNewUserInServer) {
	for user := range s.users {
		UsersEngine.Send(user, msg)
	}
}

// USERS

func (s *server) Connect(ctx *actor.Context, msg *protoTypes.Connect) {
	sender := ctx.Sender()
	if _, ok := s.users[sender]; ok {
		s.logger.Warn("user already connected to this server", "user", ctx.Sender().GetID())
		return
	}
	s.users[sender] = true
	s.logger.Info("user connected to this server", "user", ctx.Sender().GetID(), "id", ctx.PID())

	userId := utils.GetEntityIdFromPID(sender)
	serverId := utils.GetEntityIdFromPID(ctx.PID())

	for user := range s.users {
		if user == sender {
			UsersEngine.Send(user, &protoTypes.BroadcastConnect{
				ServerId: serverId,
				UserId:   userId,
				Users:    s.usersSlice,
			})
		} else {
			UsersEngine.Send(user, &protoTypes.BroadcastConnect{
				ServerId: serverId,
				UserId:   userId,
				Type:     msg.Type,
			})
		}
	}

	s.usersSlice = append(s.usersSlice, userId)

	if msg.Type == "JOIN_SERVER" {
		for _, channel := range ctx.Children() {
			ServersEngine.SendWithSender(channel, &protoTypes.Connect{}, sender)
		}
	}
}

func (s *server) Disconnect(ctx *actor.Context, msg *protoTypes.Disconnect) {
	sender := ctx.Sender()
	_, ok := s.users[sender]
	if !ok {
		s.logger.Warn("unknown user disconnected", "user", sender, "id", ctx.PID())
		return
	}
	s.logger.Info("user disconnected", "sender", ctx.Sender(), "id", ctx.PID())

	userId := utils.GetEntityIdFromPID(sender)
	serverId := utils.GetEntityIdFromPID(ctx.PID())

	idx := slices.Index(s.usersSlice, userId)
	s.usersSlice = slices.Delete(s.usersSlice, idx, idx+1)
	delete(s.users, sender)

	for user := range s.users {
		UsersEngine.Send(user, &protoTypes.BroadcastDisconnect{
			ServerId: serverId,
			UserId:   userId,
			Type:     msg.Type,
		})
	}

	if msg.Type == "LEAVE_SERVER" {
		for _, channel := range ctx.Children() {
			ServersEngine.SendWithSender(channel, &protoTypes.Disconnect{}, sender)
		}
	}
}

// CHANNELS

func (s *server) InitializeChannels(serverId string, ctx *actor.Context) {
	strSplit := strings.Split(serverId, "/")
	id := strSplit[len(strSplit)-1]

	channels, err := db.Query.GetChannelsFromServer(context.TODO(), id)
	if err != nil {
		panic(err)
	}

	for _, channel := range channels {
		actorPid := ctx.SpawnChild(NewChannel, "channel", actor.WithID(channel.ID))
		s.channels[actorPid] = true
	}
}

func (s *server) StartChannel(ctx *actor.Context, msg *protoTypes.StartChannel) {
	channelPid := ctx.SpawnChild(NewChannel, "channel", actor.WithID(msg.ChannelId))

	for _, user := range msg.Users {
		userPID := UsersEngine.Registry.GetPID("user", user)
		UsersEngine.Send(userPID, &protoTypes.ChannelStarting{
			ActorId:      channelPid.ID,
			ActorAddress: channelPid.Address,
		})
	}
}

func (s *server) CreateChannel(ctx *actor.Context, msg *protoTypes.BodyChannelCreation) {
	channelToCreate := &services.CreateChannelBody{
		Name:        msg.Name,
		Type:        db.ChannelType(msg.Type),
		Description: msg.Description,
		Users:       msg.Users,
		Roles:       msg.Roles,
		X:           msg.X,
		Y:           msg.Y,
	}

	if msg.Id != "" {
		channelToCreate.Id = &msg.Id
	}

	channel, err := services.CreateChannel(context.TODO(), msg.CreatorId, msg.ServerId, channelToCreate)
	if err != nil {
		slog.Error("failed to create channel", "err", err)
		return
	}
	channelPid := ctx.SpawnChild(NewChannel, "channel", actor.WithID(channel.Id))
	channel.ActorId = channelPid.ID
	channel.ActorAddress = channelPid.Address

	if len(msg.Users) > 0 {
		for _, user := range msg.Users {
			userPID := UsersEngine.Registry.GetPID("user", user)
			UsersEngine.Send(userPID, channel)
		}
	}

	if len(msg.Users) <= 0 && len(msg.Roles) <= 0 {
		for user := range s.users {
			UsersEngine.Send(user, channel)
		}
	}
}

func (s *server) RemoveChannel(ctx *actor.Context, msg *protoTypes.BodyChannelRemoved) {
	err := services.DeleteChannel(context.TODO(), msg.ServerId, msg.ChannelId, msg.UserId)
	if err != nil {
		slog.Error("failed to delete channel", "err", err)
		return
	}

	channelId := fmt.Sprintf("channel/%s", msg.ChannelId)
	channelPID := ctx.PID().Child(channelId)
	ctx.Engine().Poison(channelPID)
	delete(s.channels, channelPID)

	for user := range s.users {
		UsersEngine.Send(user, &protoTypes.BroadcastChannelRemoved{
			ServerId:     msg.ServerId,
			ChannelId:    msg.ChannelId,
			ActorId:      channelPID.ID,
			ActorAddress: channelPID.Address,
		})
	}
}

func (s *server) KillChannel(ctx *actor.Context, msg *protoTypes.KillChannel) {
	channelId := fmt.Sprintf("channel/%s", msg.ChannelId)
	channelPID := ctx.PID().Child(channelId)
	ctx.Engine().Poison(channelPID)
	delete(s.channels, channelPID)

	msg.ActorAddress = channelPID.Address
	msg.ActorId = channelPID.ID

	for _, user := range msg.Users {
		userPID := UsersEngine.Registry.GetPID("user", user)
		UsersEngine.Send(userPID, msg)
	}
}
