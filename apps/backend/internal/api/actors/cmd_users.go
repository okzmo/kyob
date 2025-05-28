package actors

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/lxzan/gws"
	"github.com/okzmo/kyob/db"
	protoTypes "github.com/okzmo/kyob/types"
	"google.golang.org/protobuf/proto"
)

func (u *user) InitializeUser(ctx *actor.Context) {
	strSplit := strings.Split(ctx.PID().GetID(), "/")
	userId := strSplit[len(strSplit)-1]

	servers, err := db.Query.GetServersFromUser(context.TODO(), userId)
	if err != nil {
		u.logger.Error("no servers found for the user with id", "id", userId, "err", err)
	}

	for _, server := range servers {
		serverPID := ServersEngine.Registry.GetPID("server", server.ID)
		u.servers[serverPID] = true

		ServersEngine.SendWithSender(serverPID, &protoTypes.Connect{Type: "CONNECTING"}, ctx.PID())
		channels, _ := db.Query.GetChannelsFromServer(context.TODO(), server.ID)

		for _, channel := range channels {
			channelPID := ServersEngine.Registry.GetPID(fmt.Sprintf("server/%s/channel", server.ID), channel.ID)
			u.channels[channelPID] = true

			ServersEngine.SendWithSender(channelPID, &protoTypes.Connect{Type: "CONNECTING"}, ctx.PID())
		}
	}
}

func (u *user) KillUser(ctx *actor.Context) {
	for server := range u.servers {
		ServersEngine.SendWithSender(server, &protoTypes.Disconnect{
			Type: "DISCONNECTING",
		}, ctx.PID())
	}

	for channel := range u.channels {
		ServersEngine.SendWithSender(channel, &protoTypes.Disconnect{
			Type: "DISCONNECTING",
		}, ctx.PID())
	}
}

func (u *user) NewServer(ctx *actor.Context, msg *protoTypes.NewServerCreated) {
	serverPid := actor.NewPID(msg.ActorAddress, msg.ActorId)
	u.servers[serverPid] = true
	ServersEngine.SendWithSender(serverPid, &protoTypes.Connect{Type: "CONNECTING"}, ctx.PID())
}

func (u *user) BroadcastNewUserInServer(ctx *actor.Context, msg *protoTypes.BodyNewUserInServer) {
	msgToSend := &protoTypes.WSMessage{
		Content: &protoTypes.WSMessage_NewUser{
			NewUser: &protoTypes.BroadcastNewUserInServer{
				ServerId: msg.ServerId,
				User:     msg.User,
			},
		},
	}
	m, _ := proto.Marshal(msgToSend)
	u.wsConn.WriteMessage(gws.OpcodeBinary, m)
}

func (u *user) BroadcastConnect(ctx *actor.Context, msg *protoTypes.BroadcastConnect) {
	msgToSend := &protoTypes.WSMessage{
		Content: &protoTypes.WSMessage_UserConnect{
			UserConnect: msg,
		},
	}
	m, _ := proto.Marshal(msgToSend)
	u.wsConn.WriteMessage(gws.OpcodeBinary, m)
}

func (u *user) BroadcastDisconnect(ctx *actor.Context, msg *protoTypes.BroadcastDisconnect) {
	msgToSend := &protoTypes.WSMessage{
		Content: &protoTypes.WSMessage_UserDisconnect{
			UserDisconnect: msg,
		},
	}
	m, _ := proto.Marshal(msgToSend)
	u.wsConn.WriteMessage(gws.OpcodeBinary, m)
}

func (u *user) BroadcastChannelCreation(ctx *actor.Context, msg *protoTypes.BroadcastChannelCreation) {
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
	ServersEngine.SendWithSender(channelPid, &protoTypes.Connect{Type: "CONNECTING"}, ctx.PID())
	u.channels[channelPid] = true

	m, _ := proto.Marshal(msgToSend)
	u.wsConn.WriteMessage(gws.OpcodeBinary, m)
}

func (u *user) BroadcastChannelRemoved(ctx *actor.Context, msg *protoTypes.BroadcastChannelRemoved) {
	msgToSend := &protoTypes.WSMessage{
		Content: &protoTypes.WSMessage_ChannelRemoved{
			ChannelRemoved: &protoTypes.BroadcastChannelRemoved{
				ServerId:  msg.ServerId,
				ChannelId: msg.ChannelId,
			},
		},
	}

	channelPid := actor.NewPID(msg.ActorAddress, msg.ActorId)
	ServersEngine.SendWithSender(channelPid, &protoTypes.Disconnect{Type: "DISCONNECTING"}, ctx.PID())
	m, _ := proto.Marshal(msgToSend)
	u.wsConn.WriteMessage(gws.OpcodeBinary, m)
	delete(u.channels, channelPid)
}

func (u *user) BroadcastChatMessage(ctx *actor.Context, msg *protoTypes.BroadcastChatMessage) {
	msgToSend := &protoTypes.WSMessage{
		Content: &protoTypes.WSMessage_ChatMessage{
			ChatMessage: msg,
		},
	}

	m, _ := proto.Marshal(msgToSend)
	u.wsConn.WriteMessage(gws.OpcodeBinary, m)
}

func (u *user) BroadcastEditMessage(ctx *actor.Context, msg *protoTypes.BroadcastEditMessage) {
	msgToSend := &protoTypes.WSMessage{
		Content: &protoTypes.WSMessage_EditMessage{
			EditMessage: msg,
		},
	}

	m, _ := proto.Marshal(msgToSend)
	u.wsConn.WriteMessage(gws.OpcodeBinary, m)
}

func (u *user) BroadcastDeleteMessage(ctx *actor.Context, msg *protoTypes.DeleteChatMessage) {
	msgToSend := &protoTypes.WSMessage{
		Content: &protoTypes.WSMessage_DeleteMessage{
			DeleteMessage: &protoTypes.BroadcastDeleteChatMessage{
				MessageId: msg.MessageId,
				ServerId:  msg.ServerId,
				ChannelId: msg.ChannelId,
			},
		},
	}

	m, _ := proto.Marshal(msgToSend)
	u.wsConn.WriteMessage(gws.OpcodeBinary, m)
}

func (u *user) FriendInvite(ctx *actor.Context, msg *protoTypes.SendFriendInvite) {
	msgToSend := &protoTypes.WSMessage{
		Content: &protoTypes.WSMessage_FriendInvite{
			FriendInvite: msg,
		},
	}

	m, _ := proto.Marshal(msgToSend)
	u.wsConn.WriteMessage(gws.OpcodeBinary, m)
}

func (u *user) AcceptFriend(ctx *actor.Context, msg *protoTypes.AcceptFriendInvite) {
	msgToSend := &protoTypes.WSMessage{
		Content: &protoTypes.WSMessage_AcceptFriend{
			AcceptFriend: msg,
		},
	}

	m, _ := proto.Marshal(msgToSend)
	u.wsConn.WriteMessage(gws.OpcodeBinary, m)
}

func (u *user) DeleteFriend(ctx *actor.Context, msg *protoTypes.DeleteFriend) {
	msgToSend := &protoTypes.WSMessage{
		Content: &protoTypes.WSMessage_DeleteFriend{
			DeleteFriend: msg,
		},
	}

	m, _ := proto.Marshal(msgToSend)
	u.wsConn.WriteMessage(gws.OpcodeBinary, m)
}
