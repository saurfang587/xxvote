package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"xxvote/app/model"
	"xxvote/app/tools"
)

func Index(context *gin.Context) {
	context.HTML(http.StatusOK, "index.tmpl", nil)
}

func GetVotes(context *gin.Context) {
	ret := model.GetVotes()
	context.JSON(http.StatusOK, tools.ECode{
		Data: ret,
	})
}

func GetVoteInfo(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	ret := model.GetVote(id)
	//log.Printf("[printf]ret:%+v\n", ret)
	//log.Panicf("[Panicf]ret:%+v\n", ret)
	//log.Fatalf("[Fatalf]ret:%+v", ret)
	//log.SetLevel(log.DebugLevel)
	//log.Infof("[info]ret:%+v", ret)
	//log.Errorf("[error]ret:%+v", ret)

	tools.Logger.Infoln(fmt.Sprintf("ret:%+v", ret))

	context.JSON(http.StatusOK, tools.ECode{
		Data: ret,
	})
}

func DoVote(context *gin.Context) {
	userIDStr, _ := context.Cookie("Id")
	voteIdStr, _ := context.GetPostForm("vote_id")
	optStr, _ := context.GetPostFormArray("opt[]")

	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	voteId, _ := strconv.ParseInt(voteIdStr, 10, 64)
	opt := make([]int64, 0)
	for _, v := range optStr {
		optId, _ := strconv.ParseInt(v, 10, 64)
		opt = append(opt, optId)
	}

	//查询是否投过票了
	voteUser := model.GetVoteHistory(userID, voteId)
	if len(voteUser) > 0 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10010,
			Message: "您已投过票了",
		})
		return
	}
	model.DoVote(userID, voteId, opt)
	context.JSON(http.StatusOK, tools.ECode{
		Message: "投票完成",
	})
}
