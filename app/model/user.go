package model

import (
	"fmt"
	"xxvote/app/tools"
)

func GetUser(name string) *User {
	var ret User
	if err := Conn.Table("user").Where("name = ?", name).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return &ret
}

func GetUserV1(name string) (*User, error) {
	var ret User
	err := Conn.Raw("select * from user where name = ? limit 1", name).Scan(&ret).Error
	if err != nil {
		tools.Logger.Printf("err:%s", err.Error())
		return &ret, err
	}
	return &ret, nil
}

// CreateUser 思考下这里为什么传了个指针，之前有说过。创建数据实际用到的地方并不多，因此不必转化
func CreateUser(user *User) error {
	if err := Conn.Create(user).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		return err
	}
	return nil
}
