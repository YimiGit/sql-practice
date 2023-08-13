package utils

import (
	"github.com/bwmarrin/snowflake"
)

func DistributedID(node int64) (int64, error) {
	//分布式节点
	realNode, err := snowflake.NewNode(node)
	if err != nil {
		return 0, err
	}
	return realNode.Generate().Int64(), nil
}
