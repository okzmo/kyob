package utils

import (
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"
)

var Node *snowflake.Node

func SetupSnowflake() {
	nodeIDStr := os.Getenv("NODE_ID")
	nodeID, err := strconv.Atoi(nodeIDStr)
	if err != nil {
		panic(err)
	}

	node, err := snowflake.NewNode(int64(nodeID))
	if err != nil {
		panic(err)
	}

	Node = node
}
