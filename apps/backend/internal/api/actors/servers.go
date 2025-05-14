package actors

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/okzmo/kyob/db"
	services "github.com/okzmo/kyob/internal/service"
	protoTypes "github.com/okzmo/kyob/types"
)

var ServersEngine *actor.Engine

type (
	ChannelMap map[*actor.PID]bool
	UserMap    map[*actor.PID]bool
)

type server struct {
	channels ChannelMap
	users    UserMap
	logger   *slog.Logger
}

type channel struct {
	users  UserMap
	logger *slog.Logger
}

func NewServer() actor.Receiver {
	return &server{
		channels: make(ChannelMap),
		users:    make(UserMap),
		logger:   slog.Default(),
	}
}

func NewChannel() actor.Receiver {
	return &channel{
		users:  make(UserMap),
		logger: slog.Default(),
	}
}

func (s *server) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		slog.Info("server started",
			"id", ctx.PID().GetID(),
		)
		initializeChannels(ctx.PID().GetID(), ctx, s)
	case actor.InternalError:
		slog.Error("server erroring",
			"id", ctx.PID().GetID(),
			"err", msg.Err,
		)
	case *protoTypes.Connect:
		sender := ctx.Sender()
		if _, ok := s.users[sender]; ok {
			s.logger.Warn("user already connected to this server", "user", ctx.Sender().GetID())
			return
		}
		s.users[sender] = true
		s.logger.Info("user connected to this server", "user", ctx.Sender().GetID(), "id", ctx.PID())
	case *protoTypes.Disconnect:
		sender := ctx.Sender()
		_, ok := s.users[sender]
		if !ok {
			s.logger.Warn("unknown user disconnected", "user", sender, "id", ctx.PID())
			return
		}
		s.logger.Info("user disconnected", "sender", ctx.Sender())
		delete(s.users, sender)
	case *protoTypes.BodyChannelCreation:
		channelToCreate := &services.CreateChannelBody{
			Name:        msg.Name,
			Type:        db.ChannelType(msg.Type),
			Description: msg.Description,
			X:           msg.X,
			Y:           msg.Y,
		}

		channel, err := services.CreateChannel(context.TODO(), msg.CreatorId, msg.ServerId, channelToCreate)
		if err != nil {
			slog.Error("failed to create channel", "err", err)
			return
		}
		channelPid := ctx.SpawnChild(NewChannel, "channel", actor.WithID(strconv.Itoa(int(channel.Id))))
		channel.ActorId = channelPid.ID
		channel.ActorAddress = channelPid.Address

		for user := range s.users {
			UsersEngine.Send(user, channel)
		}
	case *protoTypes.BodyChannelRemoved:
		err := services.DeleteChannel(context.TODO(), int(msg.ServerId), int(msg.ChannelId), msg.UserId)
		if err != nil {
			slog.Error("failed to delete channel", "err", err)
			return
		}

		channelId := fmt.Sprintf("channel/%s", strconv.Itoa(int(msg.ChannelId)))
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
}

func (c *channel) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.InternalError:
		slog.Info("channel erroring",
			"id", ctx.PID().GetID(),
			"err", msg.Err,
		)
	case *protoTypes.Connect:
		sender := ctx.Sender()
		if _, ok := c.users[sender]; ok {
			c.logger.Warn("user already connected", "user", ctx.Sender().GetID())
			return
		}
		c.users[sender] = true
		c.logger.Info("user connected", "sender", ctx.Sender().GetID(), "id", ctx.PID())
	case *protoTypes.Disconnect:
		sender := ctx.Sender()
		_, ok := c.users[sender]
		if !ok {
			c.logger.Warn("unknown user disconnected", "user", sender)
			return
		}
		c.logger.Info("user disconnected", "sender", ctx.Sender())
		delete(c.users, sender)
	case *protoTypes.IncomingChatMessage:
		messageToSend := &services.CreateMessageBody{
			Content: msg.Content,
		}

		message, err := services.CreateMessage(context.TODO(), msg.Author, msg.ServerId, msg.ChannelId, messageToSend)
		if err != nil {
			slog.Error("failed to create message", "err", err)
			return
		}

		for user := range c.users {
			UsersEngine.Send(user, message)
		}
	}
}

func SetupServersEngine() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}
	ServersEngine = e

	servers, err := db.Query.GetServers(context.TODO())
	if err != nil {
		panic(err)
	}

	for _, server := range servers {
		e.Spawn(NewServer, "server", actor.WithID(strconv.Itoa(int(server.ID))))
	}
}

func initializeChannels(serverId string, ctx *actor.Context, serverActor *server) {
	strSplit := strings.Split(serverId, "/")
	idStr := strSplit[len(strSplit)-1]

	id, _ := strconv.Atoi(idStr)

	channels, err := db.Query.GetChannelsFromServer(context.TODO(), int64(id))
	if err != nil {
		panic(err)
	}

	for _, channel := range channels {
		actorPid := ctx.SpawnChild(NewChannel, "channel", actor.WithID(strconv.Itoa(int(channel.ID))))
		serverActor.channels[actorPid] = true
	}
}
