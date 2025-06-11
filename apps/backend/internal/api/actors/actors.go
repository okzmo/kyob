package actors

import (
	"context"
	"log/slog"

	"github.com/anthdm/hollywood/actor"
	"github.com/lxzan/gws"
	"github.com/okzmo/kyob/db"
	protoTypes "github.com/okzmo/kyob/types"
)

var (
	UsersEngine   *actor.Engine
	ServersEngine *actor.Engine
)

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
		e.Spawn(NewServer, "server", actor.WithID(server.ID))
	}
}

type VoiceUser struct {
	Id     string `json:"user_id"`
	Deafen bool   `json:"deafen"`
	Mute   bool   `json:"mute"`
}

type (
	ChannelMap map[*actor.PID]bool
	ServerMap  map[*actor.PID]bool
	UserMap    map[*actor.PID]bool
	CallMap    map[string]VoiceUser
)

type server struct {
	channels   ChannelMap
	users      UserMap
	usersSlice []string
	logger     *slog.Logger
}

func NewServer() actor.Receiver {
	return &server{
		channels: make(ChannelMap),
		users:    make(UserMap),
		logger:   slog.Default(),
	}
}

func (s *server) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Stopped:
		slog.Info("server stopped",
			"id", ctx.PID().GetID(),
		)
	case actor.Started:
		slog.Info("server started",
			"id", ctx.PID().GetID(),
		)
		s.InitializeChannels(ctx.PID().GetID(), ctx)
	case actor.InternalError:
		slog.Error("server erroring",
			"id", ctx.PID().GetID(),
			"err", msg.Err,
		)
	case *protoTypes.Connect:
		s.Connect(ctx, msg)
	case *protoTypes.Disconnect:
		s.Disconnect(ctx, msg)
	case *protoTypes.StartChannel:
		s.StartDMChannel(ctx, msg)
	case *protoTypes.BodyChannelCreation:
		s.CreateChannel(ctx, msg)
	case *protoTypes.KillChannel:
		s.KillChannel(ctx, msg)
	case *protoTypes.BodyChannelRemoved:
		s.RemoveChannel(ctx, msg)
	case *protoTypes.BodyServerRemoved:
		s.RemoveServer(ctx, msg)
	case *protoTypes.BodyNewUserInServer:
		s.NewUser(ctx, msg)
	case *protoTypes.BroadcastUserInformations:
		s.BroadcastUserInformations(ctx, msg)
	}
}

type channel struct {
	users  UserMap
	call   CallMap
	logger *slog.Logger
}

func NewChannel() actor.Receiver {
	return &channel{
		users:  make(UserMap),
		call:   make(CallMap),
		logger: slog.Default(),
	}
}

func (c *channel) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Stopped:
		slog.Info("channel stopped",
			"id", ctx.PID().GetID(),
		)
	case actor.Started:
		slog.Info("channel started",
			"id", ctx.PID().GetID(),
		)
	case actor.InternalError:
		slog.Info("channel erroring",
			"id", ctx.PID().GetID(),
			"err", msg.Err,
		)
	case *protoTypes.Connect:
		c.Connect(ctx)
	case *protoTypes.Disconnect:
		c.Disconnect(ctx)
	case *protoTypes.ConnectToCall:
		c.ConnectToCall(ctx, msg)
	case *protoTypes.DisconnectFromCall:
		c.DisconnectFromCall(ctx, msg)
	case *protoTypes.IncomingChatMessage:
		c.NewMessage(ctx, msg)
	case *protoTypes.EditChatMessage:
		c.EditMessage(ctx, msg)
	case *protoTypes.DeleteChatMessage:
		c.DeleteMessage(ctx, msg)
	case *protoTypes.BroadcastUserInformations:
		c.BroadcastUserInformations(ctx, msg)
	}
}

type user struct {
	servers  ServerMap
	channels ChannelMap
	wsConn   *gws.Conn
	logger   *slog.Logger
}

func NewUser(wsConn *gws.Conn) actor.Producer {
	return func() actor.Receiver {
		return &user{
			servers:  make(ServerMap),
			channels: make(ChannelMap),
			wsConn:   wsConn,
			logger:   slog.Default(),
		}
	}
}

func (u *user) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		u.logger.Info("user starting")
		u.InitializeUser(ctx)
	case actor.Stopped:
		u.logger.Info("user stopped")
		u.KillUser(ctx)
	case actor.InternalError:
		slog.Info("user error", "err", msg.Err)
	case *protoTypes.NewServerCreated:
		u.NewServer(ctx, msg)
	case *protoTypes.BroadcastConnect:
		u.BroadcastConnect(ctx, msg)
	case *protoTypes.BroadcastDisconnect:
		u.BroadcastDisconnect(ctx, msg)
	case *protoTypes.BroadcastChannelCreation:
		u.BroadcastChannelCreation(ctx, msg)
	case *protoTypes.BroadcastChannelRemoved:
		u.BroadcastChannelRemoved(ctx, msg)
	case *protoTypes.BroadcastChatMessage:
		u.BroadcastChatMessage(ctx, msg)
	case *protoTypes.BroadcastEditMessage:
		u.BroadcastEditMessage(ctx, msg)
	case *protoTypes.DeleteChatMessage:
		u.BroadcastDeleteMessage(ctx, msg)
	case *protoTypes.CallInitialization:
		u.SendCallInitialization(ctx, msg)
	case *protoTypes.ConnectToCall:
		u.BroadcastConnectToCall(ctx, msg)
	case *protoTypes.DisconnectFromCall:
		u.BroadcastDisconnectFromCall(ctx, msg)
	case *protoTypes.BodyNewUserInServer:
		u.BroadcastNewUserInServer(ctx, msg)
	case *protoTypes.SendFriendInvite:
		u.FriendInvite(ctx, msg)
	case *protoTypes.AcceptFriendInvite:
		u.AcceptFriend(ctx, msg)
	case *protoTypes.DeleteFriend:
		u.DeleteFriend(ctx, msg)
	case *protoTypes.KillChannel:
		u.ChannelKilled(ctx, msg)
	case *protoTypes.UserChangedInformations:
		u.ChangingUserInformations(ctx, msg)
	case *protoTypes.BroadcastUserInformations:
		u.BroadcastUserInformations(ctx, msg)
	}
}

func SetupUsersEngine() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	UsersEngine = e
}
