package routesJoinAudit

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	"studentGrow/pkg/response"
	"studentGrow/service/JoinAudit"
	token2 "studentGrow/utils/token"
)

func ActivityRulerList(c *gin.Context) {
	token := token2.NewToken(c)
	_, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
		return
	}
	var cr mysql.Pagination
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "json解析失败")
		return
	}
	cr.Label = "ActivityRulerList"
	var ResAllMsgList = make([]JoinAudit.ResList, 0)
	ResAllMsgList, err = JoinAudit.ResListWithJSON(cr)
	if err != nil {
		response.ResponseSuccessWithMsg(c, err.Error(), []struct{}{})
		return
	}
	if !cr.All {
		_, msg, _ := mysql.OpenActivityStates()
		response.ResponseSuccessWithMsg(c, msg, ResAllMsgList)
		return
	}
	response.ResponseSuccess(c, ResAllMsgList)
}
func ActivityRulerManager(c *gin.Context) {
	token := token2.NewToken(c)
	_, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
		return
	}
	var cr JoinAudit.RecList
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "json数据解析失败")
		return
	}
	resList := JoinAudit.IsPassWithJSON(cr, "ruler_is_pass")
	response.ResponseSuccess(c, resList)
}
