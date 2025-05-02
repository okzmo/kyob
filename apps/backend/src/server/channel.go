package server

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
)

type channel struct {
	users []*actor.PID
}

func (c *channel) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case *Message:
		fmt.Println("work", msg.Data)
	}
}

func newChannel() actor.Receiver {
	return &channel{}
}
