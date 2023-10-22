package schedule

import (
	"log"
	"time"
	"xxvote/app/model"
)

func Start() {
	go voteEnd()

	return
}

func voteEnd() {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			//fmt.Printf("定时器 voteEnd 启动")
			log.Println("定时器 voteEnd 启动")
			_ = model.EndVote()
			log.Println("定时器 voteEnd 结束")
		}
	}
}
