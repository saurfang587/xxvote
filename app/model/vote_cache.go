package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func GetVoteCache(c context.Context, id int64) VoteWithOpt {
	var ret VoteWithOpt
	key := fmt.Sprintf("key_%d", id)
	fmt.Printf("key:%s\n", key)
	voteStr, err := Rdb.Get(c, key).Result()
	if err == nil || len(voteStr) > 0 {
		//存在数据
		_ = json.Unmarshal([]byte(voteStr), &ret)
		return ret
	}
	fmt.Printf("err:%s\n", err.Error())
	vote := GetVote(id)
	if vote.Vote.Id > 0 {
		//写入缓存
		s, _ := json.Marshal(vote)
		err1 := Rdb.Set(c, key, s, 3600*time.Second).Err()
		if err1 != nil {
			fmt.Printf("err1:%s\n", err1.Error())
		}
		ret = vote
	}

	return ret
}

func GetVoteHistoryV1(c context.Context, userId, voteId int64) []VoteOptUser {
	ret := make([]VoteOptUser, 0)
	//先查询缓存
	k := fmt.Sprintf("vote-user-%d", userId)
	str, _ := Rdb.Get(c, k).Result()
	fmt.Printf("str:%s\n", str)
	if len(str) > 0 {
		//将数据转化为struct
		_ = json.Unmarshal([]byte(str), &ret)
		return ret
	}

	//不存在就先查数据库再封装缓存
	if err := Conn.Table("vote_opt_user").Where("user_id = ? and vote_id = ?", userId, voteId).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}

	if len(ret) > 0 {
		s, _ := json.Marshal(ret)
		err := Rdb.Set(c, k, s, 3600*time.Second).Err()
		if err != nil {
			fmt.Printf("err1:%s\n", err.Error())
		}
	}

	return ret
}
