package model

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
	"xxvote/app/tools"
)

func GetVotes() []Vote {
	ret := make([]Vote, 0)
	if err := Conn.Table("vote").Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

func GetVotesV1() []Vote {
	ret := make([]Vote, 0)
	err := Conn.Raw("select * from vote").Scan(&ret).Error
	if err != nil {
		tools.Logger.Printf("[GetVotesV1]err:%s", err.Error())
	}
	return ret
}

func GetVote(id int64) VoteWithOpt {
	var ret Vote
	if err := Conn.Table("vote").Where("id = ?", id).First(&ret).Error; err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}

	opt := make([]VoteOpt, 0)
	if err := Conn.Table("vote_opt").Where("vote_id = ?", id).Find(&opt).Error; err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	return VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}
}

// GetVoteV1 改为原生SQL 用的较多的方式
func GetVoteV1(id int64) (*VoteWithOpt, error) {
	var ret Vote
	opt := make([]VoteOpt, 0)
	err := Conn.Raw("select * from vote where id = ?", id).Scan(&ret).Error
	if err != nil {
		tools.Logger.Printf("[GetVoteV1]err:%s", err.Error())
		return nil, err
	}

	err1 := Conn.Raw("select * from vote_opt where vote_id = ?").Scan(&opt).Error
	if err1 != nil {
		tools.Logger.Printf("[GetVoteV1]err:%s", err.Error())
		return nil, err
	}

	return &VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}, nil
}

// GetVoteV2 改为Gorm预加载模式 建议使用
func GetVoteV2(id int64) (*VoteWithOptV1, error) {
	var ret VoteWithOptV1
	err := Conn.Preload("vote_opt", "vote_id = ?", id).Raw("select * from vote where id = ?", id).Scan(&ret).Error
	if err != nil {
		tools.Logger.Printf("[GetVoteV1]err:%s", err.Error())
		return nil, err
	}

	return &ret, nil
}

// GetVoteV3 改为Join模式 数据量少的时候会用的方式
func GetVoteV3(id int64) (*VoteWithOptV1, error) {
	var ret VoteWithOptV1
	sql := "select vote.*,vote_opt.id as vid, vote_opt.name,vote_opt.count from vote join vote_opt on vote.id = vote_opt.vote_id where vote.id = ?"
	//err := Conn.Raw(sql, id).Scan(&ret).Error //这样子是无法直接扫到结构体里的,我们可以提供两个方法：
	//第一个 把ret 换成map
	//ret1 := make(map[any]any)
	//err := Conn.Raw(sql, id).Scan(&ret1).Error
	//for a, a2 := range ret1 {
	//	//再把 a a2 转义到 VoteWithOpt中
	//}
	//第二种方法
	rows, err := Conn.Raw(sql, id).Rows()
	if err != nil {
		return &ret, err
	}
	//opt := make([]VoteOpt, 0)
	for rows.Next() {
		//读取vote_opt数据
		//ret1 := make(map[string]interface{})
		_ = Conn.ScanRows(rows, &ret)

		//再将map 的数据转存到结构体中，注意，这个方法非常难用，非常不好用。
		//Gorm 还提供了一种自定义数据结构的方法，也不太好用

		fmt.Printf("ret1:%+v\n", ret)
	}

	return &ret, nil
}

// GetVoteV4 改为并发模式1 绝对不会用的方式
func GetVoteV4(id int64) (*VoteWithOpt, error) {
	var ret Vote
	opt := make([]VoteOpt, 0)

	ch := make(chan struct{}, 2)
	go func() {
		err := Conn.Raw("select * from vote where id = ?", id).Scan(&ret).Error
		if err != nil {
			tools.Logger.Printf("[GetVoteV1]err:%s", err.Error())
		}
		ch <- struct{}{}
	}()

	go func() {
		err := Conn.Raw("select * from vote_opt where vote_id = ?", id).Scan(&opt).Error
		if err != nil {
			tools.Logger.Printf("[GetVoteV1]err:%s", err.Error())
		}
		ch <- struct{}{}
	}()

	var ini int
	for _ = range ch {
		ini++
		if ini >= 2 {
			break
		}
	}

	return &VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}, nil
}

// GetVoteV5 改为并发模式2 最常用的方式
func GetVoteV5(id int64) (*VoteWithOpt, error) {
	var ret Vote
	opt := make([]VoteOpt, 0)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := Conn.Raw("select * from vote where id = ?", id).Scan(&ret).Error
		if err != nil {
			tools.Logger.Printf("[GetVoteV1]err:%s", err.Error())
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := Conn.Raw("select * from vote_opt where vote_id = ?", id).Scan(&opt).Error
		if err != nil {
			tools.Logger.Printf("[GetVoteV1]err:%s", err.Error())
		}
		wg.Done()
	}()
	wg.Wait()

	return &VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}, nil
}

// DoVote 通用方式
func DoVote(userId, voteId int64, optIDs []int64) bool {

	tx := Conn.Begin()

	var ret Vote
	//有没有这个投票
	if err := tx.Table("vote").Where("id = ?", voteId).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		tx.Rollback()
		return false
	}
	//检查是否投过票
	var oldUser VoteOptUser
	if err := tx.Table("vote_opt_user").Where("user_id = ? and vote_id = ?", userId, voteId).First(&oldUser).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		tx.Rollback()
		return false
	}
	if oldUser.Id > 0 {
		fmt.Printf("err:%s", "用户已经投过票了！")
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

// DoVoteV3 将SQL优化为原生SQL
func DoVoteV3(userId, voteId int64, optIDs []int64) bool {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		var ret Vote
		err := tx.Raw("select * from vote where id = ? limit 1", voteId).Save(&ret).Error
		if err != nil {
			tools.Logger.Printf("[DoVoteV3]err:%s", err.Error())
			return err //只要返回了err GORM会直接回滚，不会提交
		}

		for _, value := range optIDs {
			err1 := tx.Exec("update vote_opt set count = count+1 where id = ? limit 1", value).Error
			if err1 != nil {
				tools.Logger.Printf("[DoVoteV3]err1:%s", err.Error())
				return err
			}
			user := VoteOptUser{
				VoteId:      voteId,
				UserId:      userId,
				VoteOptId:   value,
				CreatedTime: time.Now(),
			}
			err2 := tx.Create(&user).Error // 通过数据的指针来创建
			if err2 != nil {
				tools.Logger.Printf("[DoVoteV3]err2:%s", err.Error())
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

func AddVote(vote Vote, opt []VoteOpt) error {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&vote).Error; err != nil {
			return err
		}
		for _, voteOpt := range opt {
			voteOpt.VoteId = vote.Id
			if err := tx.Create(&voteOpt).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
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

func DelVoteV1(id int64) error {
	if err := Conn.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec("delete from vote where id = ? limit 1", id).Error
		if err != nil {
			tools.Logger.Printf("[DelVoteV1]err:%s", err.Error())
			return err
		}

		err1 := tx.Exec("delete from vote_opt where vote_id = ?", id).Error
		if err1 != nil {
			tools.Logger.Printf("[DelVoteV1]err1:%s", err.Error())
			return err
		}

		err2 := tx.Exec("delete from vote_opt_user where vote_id = ?", id).Error
		if err2 != nil {
			tools.Logger.Printf("[DelVoteV1]err2:%s", err.Error())
			return err
		}

		return nil
	}); err != nil {
		tools.Logger.Printf("[DelVoteV1]err:%s", err.Error())
		return err
	}

	return nil
}

func UpdateVote(vote Vote, opt []VoteOpt) error {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&vote).Error; err != nil {
			return err
		}
		for _, voteOpt := range opt {
			if err := tx.Save(&voteOpt).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func GetVoteHistory(userId, voteId int64) []VoteOptUser {
	ret := make([]VoteOptUser, 0)
	if err := Conn.Table("vote_opt_user").Where("user_id = ? and vote_id = ?", userId, voteId).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

func EndVote() error {
	//执行逻辑
	votes := make([]Vote, 0)
	err := Conn.Table("vote").Where("status = 1").Find(&votes).Error
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	for _, vote := range votes {
		old := vote.CreatedTime.Unix()
		if old+vote.Time < now {
			//到期了,就关闭掉
			if err1 := Conn.Table("vote").Where("id = ?", vote.Id).Update("status", 2).Error; err1 != nil {
				fmt.Printf("err:%s", err1.Error())
			}
		}
	}

	return nil
}
