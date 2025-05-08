package actors

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/okzmo/kyob/db"
)

var ServersEngine *actor.Engine

type (
	channelMap map[string]*actor.PID
	userMap    map[string]*actor.PID
)

type server struct {
	channels channelMap
	users    userMap
	logger   *slog.Logger
}

type channel struct {
	users  userMap
	logger *slog.Logger
}

func NewServer() actor.Receiver {
	return &server{
		channels: make(channelMap),
		users:    make(userMap),
		logger:   slog.Default(),
	}
}

func NewChannel() actor.Receiver {
	return &channel{
		users:  make(userMap),
		logger: slog.Default(),
	}
}

func (s *server) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		slog.Info("new server initialized",
			"id", ctx.PID().GetID(),
		)
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

func (s *channel) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		slog.Info("new channel initialized",
			"id", ctx.PID().GetID(),
		)
	case actor.Started:
		slog.Info("channel started",
			"id", ctx.PID().GetID(),
		)
	case actor.Stopped:
		slog.Info("channel stopped",
			"id", ctx.PID().GetID(),
		)
	case actor.InternalError:
		slog.Info("channel erroring",
			"id", ctx.PID().GetID(),
			"err", msg.Err,
		)
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
		serverActor.channels[actorPid.Address] = actorPid
	}
}
