package router

import (
	"fmt"
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

	{
		index := r.Group("")
		index.Use(checkUser)
		//vote
		index.GET("/index", logic.Index) //静态页面

		index.GET("/votes", logic.GetVotes)
		index.GET("/vote", logic.GetVoteInfo)
		index.POST("/vote", logic.DoVote)

		index.POST("/vote/add", logic.AddVote)
		index.POST("/vote/update", logic.UpdateVote)
		index.POST("/vote/del", logic.DelVote)

		index.GET("/result", logic.ResultInfo)
		index.GET("/result/info", logic.ResultVote)
	}

	r.GET("/", logic.Index)

	{
		//login
		r.GET("/login", logic.GetLogin)
		r.POST("/login", logic.DoLogin)
		r.GET("/logout", logic.Logout)

		//user
		r.POST("/user/create", logic.CreateUser)
	}

	//验证码
	{
		r.GET("/captcha", func(context *gin.Context) {
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
		})

		r.POST("/captcha/verify", func(context *gin.Context) {
			var param tools.CaptchaData
			if err := context.ShouldBind(&param); err != nil {
				context.JSON(http.StatusOK, tools.ParamErr)
				return
			}

			fmt.Printf("参数为：%+v", param)
			if !tools.CaptchaVerify(param) {
				context.JSON(http.StatusOK, tools.ECode{
					Code:    10008,
					Message: "验证失败",
				})
				return
			}
			context.JSON(http.StatusOK, tools.OK)
		})
	}

	if err := r.Run(":8080"); err != nil {
		panic("gin 启动失败！")
	}
}

func checkUser(context *gin.Context) {
	var name string
	var id int64 //TODO 存在一个bug
	values := model.GetSession(context)

	if v, ok := values["name"]; ok {
		name = v.(string)
	}

	if v, ok := values["id"]; ok {
		id = v.(int64)
	}

	if name == "" || id <= 0 {
		//context.JSON(http.StatusOK, tools.NotLogin)
		//context.Abort()
	}

	context.Next()
}
