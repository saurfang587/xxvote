package app

import (
	"xxvote/app/model"
	"xxvote/app/router"
)

// Start 启动器方法
func Start() {
	model.NewMysql()
	defer func() {
		model.Close()
	}()

	router.New()
}
