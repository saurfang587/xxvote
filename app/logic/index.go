package logic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"xxvote/app/model"
	"xxvote/app/tools"
)

func Index(context *gin.Context) {
	ret := model.GetVotes()
	context.HTML(http.StatusOK, "index.tmpl", gin.H{"vote": ret})
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
	//log.Printf("[print]ret:%+v\n", ret)
	//log.Panicf("[panic]ret:%+v\n", ret)
	//log.Fatalf("[fatal]ret:%+v\n", ret)
	tools.Logger.Errorf("[error]ret:%+v", ret)
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

	old := model.GetVoteHistory(userID, voteId)
	if len(old) >= 1 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10010,
			Message: "您已投过票",
		})
	}

	opt := make([]int64, 0)
	for _, v := range optStr {
		optId, _ := strconv.ParseInt(v, 10, 64)
		opt = append(opt, optId)
	}

	model.DoVote(userID, voteId, opt)
	context.JSON(http.StatusOK, tools.ECode{
		Message: "投票完成",
	})
}
