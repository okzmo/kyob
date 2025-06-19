package actors

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/anthdm/hollywood/actor"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
	protoTypes "github.com/okzmo/kyob/types"
)

// USERS

func (c *channel) Connect(ctx *actor.Context) {
	sender := ctx.Sender()
	channelID := utils.GetEntityIdFromPID(ctx.PID())
	serverID := utils.GetEntityIdFromPID(ctx.Parent())

	if _, ok := c.users[sender]; ok {
		c.logger.Warn("user already connected", "user", ctx.Sender().GetID())
		return
	}
	c.users[sender] = true
	c.logger.Info("user connected", "sender", ctx.Sender().GetID(), "id", ctx.PID())

	if len(c.call) > 0 {
		var callUsers []*protoTypes.ConnectToCall
		for _, v := range c.call {
			callUsers = append(callUsers, &protoTypes.ConnectToCall{
				UserId:    v.ID,
				ServerId:  serverID,
				ChannelId: channelID,
			})
		}

		UsersEngine.Send(sender, &protoTypes.CallInitialization{
			CallUsers: callUsers,
		})
	}
}

func (c *channel) Disconnect(ctx *actor.Context) {
	sender := ctx.Sender()
	senderID := utils.GetEntityIdFromPID(sender)
	channelID := utils.GetEntityIdFromPID(ctx.PID())
	serverID := utils.GetEntityIdFromPID(ctx.Parent())

	_, ok := c.users[sender]
	if !ok {
		c.logger.Warn("unknown user disconnected", "user", sender, "id", ctx.PID())
		return
	}
	c.logger.Info("user disconnected", "sender", ctx.Sender(), "id", ctx.PID())
	delete(c.users, sender)

	if _, ok := c.call[senderID]; ok {
		delete(c.call, senderID)

		for user := range c.users {
			UsersEngine.Send(user, &protoTypes.DisconnectFromCall{
				UserId:    senderID,
				ChannelId: channelID,
				ServerId:  serverID,
			})
		}
	}
}

// MESSAGES

func (c *channel) NewMessage(ctx *actor.Context, msg *protoTypes.IncomingChatMessage) {
	messageToSend := &services.MessageBody{
		Content:       msg.Content,
		Everyone:      msg.Everyone,
		MentionsUsers: msg.MentionsUsers,
		Attachments:   msg.Attachments,
	}

	message, err := services.CreateMessage(context.TODO(), msg.AuthorId, msg.ServerId, msg.ChannelId, messageToSend)
	if err != nil {
		slog.Error("failed to create message", "err", err)
		return
	}

	for user := range c.users {
		UsersEngine.Send(user, message)
	}
}

func (c *channel) EditMessage(ctx *actor.Context, msg *protoTypes.EditChatMessage) {
	messageToEdit := &services.MessageBody{
		Content:       msg.Content,
		Everyone:      msg.Everyone,
		MentionsUsers: msg.MentionsUsers,
	}

	message, err := services.EditMessage(context.TODO(), msg.UserId, msg.ServerId, msg.ChannelId, msg.MessageId, messageToEdit)
	if err != nil {
		slog.Error("failed to create message", "err", err)
		return
	}

	for user := range c.users {
		UsersEngine.Send(user, message)
	}
}

func (c *channel) DeleteMessage(ctx *actor.Context, msg *protoTypes.DeleteChatMessage) {
	err := services.DeleteMessage(context.TODO(), msg.MessageId, msg.UserId)
	if err != nil {
		slog.Error("failed to delete message", "err", err)
		return
	}

	for user := range c.users {
		UsersEngine.Send(user, msg)
	}
}

func (c *channel) BroadcastUserInformations(ctx *actor.Context, msg *protoTypes.BroadcastUserInformations) {
	for user := range c.users {
		UsersEngine.Send(user, msg)
	}
}

// CALL

func (c *channel) ConnectToCall(ctx *actor.Context, msg *protoTypes.ConnectToCall) {
	fmt.Println(msg.UserId)
	c.call[msg.UserId] = VoiceUser{
		ID:     msg.UserId,
		Deafen: false,
		Mute:   false,
	}

	for user := range c.users {
		UsersEngine.Send(user, msg)
	}
}

func (c *channel) DisconnectFromCall(ctx *actor.Context, msg *protoTypes.DisconnectFromCall) {
	delete(c.call, msg.UserId)

	for user := range c.users {
		UsersEngine.Send(user, msg)
	}
}
