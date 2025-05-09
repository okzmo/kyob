package actors

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/okzmo/kyob/db"
	services "github.com/okzmo/kyob/internal/service"
	proto "github.com/okzmo/kyob/types"
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
	case actor.Stopped:
		slog.Info("server stopped",
			"id", ctx.PID().GetID(),
		)
	case actor.InternalError:
		slog.Error("server erroring",
			"id", ctx.PID().GetID(),
			"err", msg.Err,
		)
	}
}

func (c *channel) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.InternalError:
		slog.Info("channel erroring",
			"id", ctx.PID().GetID(),
			"err", msg.Err,
		)
	case *proto.ConnectToChannel:
		sender := ctx.Sender()
		if _, ok := c.users[sender]; ok {
			c.logger.Warn("user already connected", "user", ctx.Sender().GetID())
			return
		}
		c.users[sender] = true
		c.logger.Info("user connected", "sender", ctx.Sender())
	case *proto.DisconnectFromChannel:
		sender := ctx.Sender()
		_, ok := c.users[sender]
		if !ok {
			c.logger.Warn("unknown user disconnected", "user", sender)
			return
		}
		c.logger.Info("user disconnected", "sender", ctx.Sender())
		delete(c.users, sender)
	case *proto.IncomingChatMessage:
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
