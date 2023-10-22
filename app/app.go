package app

import (
	"xxvote/app/model"
	"xxvote/app/router"
	"xxvote/app/tools"
)

// Start 启动器方法
func Start() {
	model.NewMysql()
	defer func() {
		model.Close()
	}()

	//schedule.Start()

	tools.NewLogger()

	router.New()
}
