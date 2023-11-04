package tools

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

func GetUUID() string {
	id := uuid.New() //默认V4 版本
	fmt.Printf("uuid:%s,version:%s", id.String(), id.Version().String())
	return id.String()
}

func GetUid() int64 {
	node, _ := snowflake.NewNode(1)
	return node.Generate().Int64()
}
