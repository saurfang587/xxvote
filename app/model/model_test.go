package model

import (
	"fmt"
	"testing"
)

func TestGetVotes(t *testing.T) {
	NewMysql()
	//测试用例
	r := GetVotes()
	fmt.Printf("ret:%+v", r)
	Close()
}

func TestGetVote(t *testing.T) {
	NewMysql()
	//测试用例
	r := GetVote(1)
	fmt.Printf("ret:%+v", r)
	Close()
}

func TestDoVote(t *testing.T) {
	NewMysql()
	//测试用例
	r := DoVote(1, 1, []int64{1, 2})
	fmt.Printf("ret:%+v", r)
	Close()
}
