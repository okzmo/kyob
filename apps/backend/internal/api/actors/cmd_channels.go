package actors

import (
	"context"
	"log/slog"

	"github.com/anthdm/hollywood/actor"
	services "github.com/okzmo/kyob/internal/service"
	protoTypes "github.com/okzmo/kyob/types"
)

// USERS

func (c *channel) Connect(ctx *actor.Context) {
	sender := ctx.Sender()
	if _, ok := c.users[sender]; ok {
		c.logger.Warn("user already connected", "user", ctx.Sender().GetID())
		return
	}
	c.users[sender] = true
	c.logger.Info("user connected", "sender", ctx.Sender().GetID(), "id", ctx.PID())
}

func (c *channel) Disconnect(ctx *actor.Context) {
	sender := ctx.Sender()
	_, ok := c.users[sender]
	if !ok {
		c.logger.Warn("unknown user disconnected", "user", sender, "id", ctx.PID())
		return
	}
	c.logger.Info("user disconnected", "sender", ctx.Sender(), "id", ctx.PID())
	delete(c.users, sender)
}

// MESSAGES

func (c *channel) NewMessage(ctx *actor.Context, msg *protoTypes.IncomingChatMessage) {
	messageToSend := &services.MessageBody{
		Content:       msg.Content,
		MentionsUsers: msg.MentionsUsers,
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

func (c *channel) EditMessage(ctx *actor.Context, msg *protoTypes.EditChatMessage) {
	messageToEdit := &services.MessageBody{
		Content:       msg.Content,
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
