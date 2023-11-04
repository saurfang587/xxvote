package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rbcervilla/redisstore/v9"
)

var sessionStore *redisstore.RedisStore
var sessionNameV1 = "session-name-v1"

func GetSessionV1(c *gin.Context) map[interface{}]interface{} {
	session, _ := sessionStore.Get(c.Request, sessionNameV1)
	fmt.Printf("session:%+v\n", session.Values)
	return session.Values
}

func SetSessionV1(c *gin.Context, name string, id int64) error {
	session, _ := sessionStore.Get(c.Request, sessionNameV1)
	session.Values["name"] = name
	session.Values["id"] = id
	return session.Save(c.Request, c.Writer)
}

func FlushSessionV1(c *gin.Context) error {
	session, _ := sessionStore.Get(c.Request, sessionNameV1)
	fmt.Printf("session : %+v\n", session.Values)
	session.Values["name"] = ""
	session.Values["id"] = ""
	return session.Save(c.Request, c.Writer)
}
