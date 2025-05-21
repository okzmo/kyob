package utils

import (
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"
)

var Node *snowflake.Node

func SetupSnowflake() {
	nodeIdStr := os.Getenv("NODE_ID")
	nodeId, err := strconv.Atoi(nodeIdStr)
	if err != nil {
		panic(err)
	}

	node, err := snowflake.NewNode(int64(nodeId))
	if err != nil {
		panic(err)
	}

	Node = node
}
