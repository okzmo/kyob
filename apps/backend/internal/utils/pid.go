package utils

import (
	"strings"

	"github.com/anthdm/hollywood/actor"
)

func GetEntityIdFromPID(entityId *actor.PID) string {
	strSplit := strings.Split(entityId.String(), "/")
	id := strSplit[len(strSplit)-1]

	return id
}
