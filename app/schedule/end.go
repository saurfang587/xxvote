package schedule

import (
	"fmt"
	"time"
	"xxvote/app/model"
)

func Start() {
	go EndVote()
}

func EndVote() {
	t := time.NewTicker(5 * time.Second)
	defer func() {
		t.Stop()
	}()

	for {
		select {
		case <-t.C:
			fmt.Println("EndVote 启动")
			//执行函数
			model.EndVote()
			fmt.Println("EndVote 运行完毕")
		}
	}

}
