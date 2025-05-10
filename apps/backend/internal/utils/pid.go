package utils

import (
	"strconv"
	"strings"

	"github.com/anthdm/hollywood/actor"
)

func GetEntityIdFromPID(entityId *actor.PID) int {
	strSplit := strings.Split(entityId.String(), "/")
	idStr := strSplit[len(strSplit)-1]
	id, _ := strconv.Atoi(idStr)

	return id
}
