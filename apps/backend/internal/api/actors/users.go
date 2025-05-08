package actors

import (
	"log/slog"

	"github.com/anthdm/hollywood/actor"
	"github.com/lxzan/gws"
)

var UsersEngine *actor.Engine

type (
	serversMap map[string]*actor.PID
)

type user struct {
	servers serversMap
	wsConn  *gws.Conn
	logger  *slog.Logger
}

func NewUser(wsConn *gws.Conn) actor.Producer {
	return func() actor.Receiver {
		return &user{
			servers: make(serversMap),
			wsConn:  wsConn,
			logger:  slog.Default(),
		}
	}
}

func (s *user) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		slog.Info("new user started", "msg", msg)
	case actor.Stopped:
		slog.Info("new user stopped", "msg", msg)
	}
}

func SetupUsersEngine() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	UsersEngine = e
}
