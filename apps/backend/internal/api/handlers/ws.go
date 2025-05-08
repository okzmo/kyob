package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-chi/chi/v5"
	"github.com/lxzan/gws"
	"github.com/okzmo/kyob/internal/api/actors"
)

const (
	PingInterval = 5 * time.Second
	PingWait     = 10 * time.Second
)

var (
	Upgrader *gws.Upgrader
	usersMap map[*gws.Conn]*actor.PID
)

func SetupWebsocket() {
	usersMap = make(map[*gws.Conn]*actor.PID)
	Upgrader = gws.NewUpgrader(&WSHandler{}, &gws.ServerOption{
		ParallelEnabled:   true,
		Recovery:          gws.Recovery,
		PermessageDeflate: gws.PermessageDeflate{Enabled: true},
	})
}

type WSHandler struct{}

func (c *WSHandler) OnOpen(socket *gws.Conn) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
}

func (c *WSHandler) OnClose(socket *gws.Conn, err error) {
	userPID := usersMap[socket]
	actors.UsersEngine.Poison(userPID)
}

func (c *WSHandler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	_ = socket.WritePong(nil)
}

func (c *WSHandler) OnPong(socket *gws.Conn, payload []byte) {}

func (c *WSHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	socket.WriteMessage(message.Opcode, message.Bytes())
}

func WS(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "user_id")

	socket, err := Upgrader.Upgrade(w, r)
	if err != nil {
		slog.Error("failed upgrading connection", "err", err)
	}

	userPID := actors.UsersEngine.Spawn(actors.NewUser(socket), "user", actor.WithID(idParam))
	usersMap[socket] = userPID

	go func() {
		socket.ReadLoop()
	}()
}
