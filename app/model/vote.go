package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

func GetVotes() []Vote {
	ret := make([]Vote, 0)
	if err := Conn.Table("vote").Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

func GetVote(id int64) VoteWithOpt {
	var ret Vote
	if err := Conn.Table("vote").Where("id = ?", id).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}

	opt := make([]VoteOpt, 0)
	if err := Conn.Table("vote_opt").Where("vote_id = ?", id).Find(&opt).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}
}

// DoVote 通用方式
func DoVote(userId, voteId int64, optIDs []int64) bool {

	tx := Conn.Begin()

	var ret Vote
	if err := tx.Table("vote").Where("id = ?", voteId).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		tx.Rollback()
		return false
	}

	for _, value := range optIDs {
		if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			tx.Rollback()
			return false
		}
		user := VoteOptUser{
			VoteId:      voteId,
			UserId:      userId,
			VoteOptId:   value,
			CreatedTime: time.Now(),
		}
		err := tx.Create(&user).Error // 通过数据的指针来创建
		if err != nil {
			fmt.Printf("err:%s", err.Error())
			tx.Rollback()
			return false
		}
	}
	tx.Commit()
	return true
}

// DoVoteV1 原生SQL
func DoVoteV1(userId, voteId int64, optIDs []int64) bool {
	Conn.Exec("begin").
		Exec("select * from vote where id = ?", voteId).
		Exec("commit")
	return false
}

// DoVoteV2 匿名函数
func DoVoteV2(userId, voteId int64, optIDs []int64) bool {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		var ret Vote
		if err := tx.Table("vote").Where("id = ?", voteId).First(&ret).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err //只要返回了err GORM会直接回滚，不会提交
		}

		for _, value := range optIDs {
			if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
				fmt.Printf("err:%s", err.Error())
				return err
			}
			user := VoteOptUser{
				VoteId:      voteId,
				UserId:      userId,
				VoteOptId:   value,
				CreatedTime: time.Now(),
			}
			err := tx.Create(&user).Error // 通过数据的指针来创建
			if err != nil {
				fmt.Printf("err:%s", err.Error())
				return err
			}
		}
		return nil //如果返回nil 则直接commit
	})

	if err != nil {
		return false
	}

	return true
}

func AddVote() {

}

func DelVote(id int64) error {
	if err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Vote{}, id).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}

		if err := tx.Where("vote_id = ?", id).Delete(&VoteOpt{}).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}

		if err := tx.Where("vote_id = ?", id).Delete(&VoteOptUser{}).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}

		return nil
	}); err != nil {
		fmt.Printf("err:%s", err.Error())
		return err
	}

	return nil
}

func UpdateVote() {

}
