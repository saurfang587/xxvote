package model

import (
	"fmt"
)

func GetUser(name string) *User {
	var ret User
	if err := Conn.Table("user").Where("name = ?", name).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return &ret
}

// CreateUser 思考下这里为什么传了个指针，之前有说过
func CreateUser(user *User) error {
	if err := Conn.Create(user).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		return err
	}
	return nil
}
