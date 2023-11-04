package app

import (
	"xxvote/app/model"
	"xxvote/app/router"
	"xxvote/app/tools"
)

// Start 启动器方法
func Start() {
	model.NewMysql()
	model.NewRdb()
	defer func() {
		model.Close()
	}()
	//schedule.Start()

	tools.NewLogs()
	router.New()
}
