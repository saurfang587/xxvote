package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
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
	model.Rdb.Set(context, fmt.Sprintf("key_%d", userID), nil, 0)
	context.JSON(http.StatusOK, tools.ECode{
		Message: "投票完成",
	})
}

func GetCaptcha(context *gin.Context) {
	//增加一个简单的限流
	if !checkXYZ(context) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: "您的手速可真快啊！",
		})
		return
	}
	captcha, err := tools.CaptchaGenerate()
	if err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, tools.ECode{
		Data: captcha,
	})
}

func checkXYZ(context *gin.Context) bool {
	//拿到IP和UA
	ip := context.ClientIP()
	ua := context.GetHeader("user-agent")
	fmt.Printf("ip:%s\nua:%s\n", ip, ua)
	//转下MD5
	hash := md5.New()
	hash.Write([]byte(ip + ua))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	//校验是否被ban
	flag, _ := model.Rdb.Get(context, "ban-"+hashString).Bool()
	if flag {
		return false
	}

	i, _ := model.Rdb.Get(context, "xyz-"+hashString).Int()
	fmt.Printf("i:%d\n", i)
	if i > 5 {
		model.Rdb.SetEx(context, "ban-"+hashString, true, 30*time.Second)
		return false
	}

	model.Rdb.Incr(context, "xyz-"+hashString)
	model.Rdb.Expire(context, "xyz-"+hashString, 50*time.Second)
	return true
}
