package actors

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/lxzan/gws"
	"github.com/okzmo/kyob/db"
	protoTypes "github.com/okzmo/kyob/types"
	"google.golang.org/protobuf/proto"
)

var UsersEngine *actor.Engine

type (
	serversMap map[string]*actor.PID
	channels   []*actor.PID
)

type user struct {
	servers  serversMap
	channels channels
	wsConn   *gws.Conn
	logger   *slog.Logger
}

func NewUser(wsConn *gws.Conn) actor.Producer {
	return func() actor.Receiver {
		return &user{
			servers:  make(serversMap),
			channels: []*actor.PID{},
			wsConn:   wsConn,
			logger:   slog.Default(),
		}
	}
}

func (u *user) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		u.logger.Info("user starting")
		u.initializeUser(ctx)
	case actor.Stopped:
		u.logger.Info("user stopped")
		u.killUser(ctx)
	case actor.InternalError:
		slog.Info("user error", "err", msg.Err)
	case *protoTypes.BroadcastChannelCreation:
		msgToSend := &protoTypes.WSMessage{
			Content: &protoTypes.WSMessage_ChannelCreation{
				ChannelCreation: &protoTypes.BroadcastChannelCreation{
					Id:          msg.Id,
					ServerId:    msg.ServerId,
					Name:        msg.Name,
					Type:        msg.Type,
					Description: msg.Description,
					Users:       msg.Users,
					Roles:       msg.Roles,
					X:           msg.X,
					Y:           msg.Y,
					CreatedAt:   msg.CreatedAt,
					UpdatedAt:   msg.UpdatedAt,
				},
			},
		}

		channelPid := actor.NewPID(msg.ActorAddress, msg.ActorId)
		ServersEngine.SendWithSender(channelPid, &protoTypes.Connect{}, ctx.PID())
		u.channels = append(u.channels, channelPid)

		m, _ := proto.Marshal(msgToSend)
		u.wsConn.WriteMessage(gws.OpcodeBinary, m)
	case *protoTypes.BroadcastChatMessage:
		msgToSend := &protoTypes.WSMessage{
			Content: &protoTypes.WSMessage_ChatMessage{
				ChatMessage: msg,
			},
		}

		m, _ := proto.Marshal(msgToSend)
		u.wsConn.WriteMessage(gws.OpcodeBinary, m)
	case *protoTypes.NewServerCreated:
		serverPid := actor.NewPID(msg.ActorAddress, msg.ActorId)
		u.servers[msg.ActorAddress] = serverPid
		ServersEngine.SendWithSender(serverPid, &protoTypes.Connect{}, ctx.PID())
	case *protoTypes.BroadcastChannelRemoved:
		msgToSend := &protoTypes.WSMessage{
			Content: &protoTypes.WSMessage_ChannelRemoved{
				ChannelRemoved: &protoTypes.BroadcastChannelRemoved{
					ServerId:  msg.ServerId,
					ChannelId: msg.ChannelId,
				},
			},
		}

		channelPid := actor.NewPID(msg.ActorAddress, msg.ActorId)
		ServersEngine.SendWithSender(channelPid, &protoTypes.Disconnect{}, ctx.PID())
		m, _ := proto.Marshal(msgToSend)
		u.wsConn.WriteMessage(gws.OpcodeBinary, m)
	}
}

func SetupUsersEngine() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	UsersEngine = e
}

func (u *user) initializeUser(ctx *actor.Context) {
	strSplit := strings.Split(ctx.PID().GetID(), "/")
	idStr := strSplit[len(strSplit)-1]
	userId, _ := strconv.Atoi(idStr)
	servers, err := db.Query.GetServersFromUser(context.TODO(), int64(userId))
	if err != nil {
		u.logger.Error("no servers found for the user with id", "id", userId, "err", err)
	}

	for _, server := range servers {
		serverPID := ServersEngine.Registry.GetPID("server", strconv.Itoa(int(server.ID)))
		u.servers[serverPID.Address] = serverPID

		ServersEngine.SendWithSender(serverPID, &protoTypes.Connect{}, ctx.PID())
		channels, _ := db.Query.GetChannelsFromServer(context.TODO(), server.ID)

		for _, channel := range channels {
			channelPID := ServersEngine.Registry.GetPID(fmt.Sprintf("server/%d/channel", server.ID), strconv.Itoa(int(channel.ID)))
			u.channels = append(u.channels, channelPID)

			ServersEngine.SendWithSender(channelPID, &protoTypes.Connect{}, ctx.PID())
		}
	}
}

func (u *user) killUser(ctx *actor.Context) {
	for _, server := range u.servers {
		ServersEngine.SendWithSender(server, &protoTypes.Disconnect{}, ctx.PID())
	}

	for _, channel := range u.channels {
		ServersEngine.SendWithSender(channel, &protoTypes.Disconnect{}, ctx.PID())
	}
}
