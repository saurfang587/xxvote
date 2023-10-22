package tools

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var Logger *logrus.Entry

func NewLogger() {
	logStore := logrus.New()
	logStore.SetLevel(logrus.DebugLevel)

	// 同时写到多个输出
	w1 := os.Stdout //写到控制台
	w2, _ := os.OpenFile("./vote.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	logStore.SetOutput(io.MultiWriter(w1, w2)) // io.MultiWriter 返回一个 io.Writer 对象

	logStore.SetFormatter(&logrus.JSONFormatter{})
	Logger = logStore.WithFields(logrus.Fields{
		"name": "香香编程喵喵喵", //在我们自己的服务里 看着很多余
		"app":  "voteV2",
	})

	//logStore.AddHook() //出现非常严重问题，直接邮箱或者微信报警，日志分割，塞入一些自定义的字段

	//context
	//logStore.WithContext()
}
