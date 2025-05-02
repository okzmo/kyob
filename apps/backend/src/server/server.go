package server

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
)

type server struct {
	channels []*actor.PID
}

type Message struct {
	Data string
}

func (s *server) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		fmt.Println("new server has initialized with pid:", ctx.PID())
	case actor.Started:
		fmt.Println("new server has started with id:", ctx.PID().ID)
		pid := ctx.SpawnChild(newChannel, "channel", actor.WithID("1"))
		fmt.Println("Spawning a channel with pid:", pid)
		s.channels = append(s.channels, pid)
	case actor.Stopped:
		fmt.Println("new server has stopped")
	case *Message:
		fmt.Println("message in server", msg.Data)
	}
}

func NewServer(id uint) actor.Producer {
	return func() actor.Receiver {
		return &server{}
	}
}
