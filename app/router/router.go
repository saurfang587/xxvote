package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xxvote/app/logic"
	"xxvote/app/model"
	"xxvote/app/tools"
)

func New() {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")
	//相关的路径 放在这里

	index := r.Group("")
	index.Use(checkUser)
	//投票相关
	{
		index.GET("/index", logic.Index)
		index.GET("/votes", logic.GetVotes)
		index.GET("/vote", logic.GetVoteInfo)
		index.POST("/vote", logic.DoVote)

		index.GET("/vote/delete", logic.DelVote)
		index.POST("/vote/add", logic.AddVote)
		index.POST("/vote/update", logic.UpdateVote)

		index.GET("/result", logic.ResultInfo)
		index.GET("/result/info", logic.ResultVote)
	}

	r.GET("/", logic.Index)

	//登录相关
	{
		r.GET("/login", logic.GetLogin)
		r.POST("/login", logic.DoLogin)
		r.GET("/logout", logic.Logout)
	}

	if err := r.Run(":8080"); err != nil {
		panic("gin 启动失败！")
	}
}

func checkUser(context *gin.Context) {
	var name string
	var id int64
	session := model.GetSession(context)
	if v, ok := session["name"]; ok {
		name = v.(string)
	}

	if v, ok := session["id"]; ok {
		id = v.(int64)
	}

	if id <= 0 || name == "" {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: "您没有登录",
		})
		context.Abort()
	}

	context.Next()
}
