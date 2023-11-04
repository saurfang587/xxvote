package model

import (
	"context"
	"fmt"
	"testing"
	"time"
	"xxvote/app/tools"
)

func TestDelVote(t *testing.T) {
	NewMysql()
	//测试用例
	r := DelVote(1)
	fmt.Printf("ret:%+v", r)
	Close()
}

func TestAddVote(t *testing.T) {
	NewMysql()
	//测试用例
	newVote := Vote{
		Title:       "测试用例",
		Type:        1,
		Status:      1,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	newVoteOpt := []VoteOpt{
		{
			Name:        "测试用例1",
			VoteId:      0,
			Count:       0,
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		},
		{
			Name:        "测试用例2",
			VoteId:      0,
			Count:       0,
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		},
	}
	r := AddVote(newVote, newVoteOpt)
	fmt.Printf("ret:%+v", r)
	Close()
}

func TestUpdateVote(t *testing.T) {
	NewMysql()
	//测试用例
	r := DelVote(1)
	fmt.Printf("ret:%+v", r)
	Close()
}

func TestGetVoteV3(t *testing.T) {
	NewMysql()
	tools.NewLogs()
	r, _ := GetVoteV5(2)
	fmt.Printf("vote:%+v", r)
}

func TestGetVoteHistoryV1(t *testing.T) {
	NewMysql()
	NewRdb()
	r := GetVoteHistoryV1(context.TODO(), 1, 1)
	fmt.Printf("vote:%+v", r)
}
