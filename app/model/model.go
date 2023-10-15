package model

import "time"

type VoteOptUser struct {
	Id          int64     `gorm:"column:id;primary_key;NOT NULL"`
	UserId      int64     `gorm:"column:user_id;default:NULL"`
	VoteId      int64     `gorm:"column:vote_id;default:NULL"`
	VoteOptId   int64     `gorm:"column:vote_opt_id;default:NULL"`
	CreatedTime time.Time `gorm:"column:created_time;default:NULL"`
}

func (v *VoteOptUser) TableName() string {
	return "vote_opt_user"
}

type VoteOpt struct {
	Id          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name        string    `gorm:"column:name;default:NULL"`
	VoteId      int64     `gorm:"column:vote_id;default:NULL"`
	Count       int64     `gorm:"column:count;default:NULL"`
	CreatedTime time.Time `gorm:"column:created_time;default:NULL"`
	UpdatedTime time.Time `gorm:"column:updated_time;default:NULL"`
}

func (v *VoteOpt) TableName() string {
	return "vote_opt"
}

type Vote struct {
	Id          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Title       string    `gorm:"column:title;default:NULL"`
	Type        int32     `gorm:"column:type;default:NULL;comment:'0单选1多选'"`
	Status      int32     `gorm:"column:status;default:NULL;comment:'0正常1超时'"`
	Time        int64     `gorm:"column:time;default:NULL;comment:'有效时长'"`
	UserId      int64     `gorm:"column:user_id;default:NULL;comment:'创建人'"`
	CreatedTime time.Time `gorm:"column:created_time;default:NULL"`
	UpdatedTime time.Time `gorm:"column:updated_time;default:NULL"`
}

func (v *Vote) TableName() string {
	return "vote"
}

type User struct {
	Id          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name        string    `gorm:"column:name;default:NULL"`
	Password    string    `gorm:"column:password;default:NULL"`
	CreatedTime time.Time `gorm:"column:created_time;default:NULL"`
	UpdatedTime time.Time `gorm:"column:updated_time;default:NULL"`
}

type VoteWithOpt struct {
	Vote Vote
	Opt  []VoteOpt
}
