package model

import (
	"fmt"
	"testing"
)

func TestDelVote(t *testing.T) {
	NewMysql()
	//测试用例
	r := DelVote(1)
	fmt.Printf("ret:%+v", r)
	Close()
}
