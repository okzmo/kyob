package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/okzmo/kyob/src/server"
)

func main() {
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	engine.Spawn(server.NewServer(1), "server", actor.WithID("1"))
	pid := engine.Registry.GetPID(fmt.Sprintf("server/%d/channel", 1), strconv.Itoa(1))
	engine.Send(pid, &server.Message{Data: "hello world"})
	time.Sleep(10 * time.Second)
}
