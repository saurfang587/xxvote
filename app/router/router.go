package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"xxvote/app/logic"
	"xxvote/app/model"
	"xxvote/app/tools"
	_ "xxvote/docs"
)

func New() {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")
	//相关的路径 放在这里

	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	index := r.Group("")
	index.Use(checkUser)
	//投票相关
	{
		index.GET("/index", logic.Index)
		index.POST("/vote", logic.DoVote)

		index.DELETE("/vote/delete", logic.DelVote)
		index.POST("/vote/add", logic.AddVote)
		index.PUT("/vote/update", logic.UpdateVote)

		index.GET("/result", logic.ResultInfo)
		index.GET("/result/info", logic.ResultVote)
	}
	//投票项目 RESTFUL 风格
	{
		index.GET("/votes", logic.GetVotes)
		index.GET("/vote", logic.GetVoteInfo)
		index.PUT("/vote", logic.UpdateVote)
		index.DELETE("/vote", logic.DelVote)

		index.POST("/do_vote", logic.DoVote)

		index.GET("/vote/result", logic.ResultVote)
	}

	r.GET("/", logic.Index)

	//登录相关
	{
		r.GET("/login", logic.GetLogin)
		r.POST("/login", logic.DoLogin)
		r.GET("/logout", logic.Logout)

		r.POST("/register", logic.CreateUser)
	}

	//验证码
	{
		r.GET("/captcha", logic.GetCaptcha)

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

	//验证 Redis
	{
		r.GET("/redis", func(context *gin.Context) {
			str := model.GetVoteCache(context, 2)
			fmt.Printf("XXBC:%+v\n", str)
		})
	}
	if err := r.Run(":8080"); err != nil {
		panic("gin 启动失败！")
	}
}

func checkUser(context *gin.Context) {
	var name string
	var id int64
	session := model.GetSessionV1(context)
	if v, ok := session["name"]; ok {
		name = v.(string)
	}

	if v, ok := session["id"]; ok {
		id = v.(int64)
	}

	if id <= 0 || name == "" {
		//context.JSON(http.StatusOK, tools.ECode{
		//	Code:    10001,
		//	Message: "您没有登录",
		//})
		//context.Abort()
	}

	context.Next()
}
