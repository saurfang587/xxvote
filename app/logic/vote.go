package logic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"xxvote/app/model"
	"xxvote/app/tools"
)

// AddVote 新增一个投票
func AddVote(context *gin.Context) {

}

// DelVote 新增一个投票
func DelVote(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	if err := model.DelVote(id); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, tools.OK)
	return
}

// UpdateVote 更新一个投票
func UpdateVote(context *gin.Context) {

}
