package routesJoinAudit

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	"studentGrow/pkg/response"
	"studentGrow/service/JoinAudit"
	token2 "studentGrow/utils/token"
)

// 组织部获取列表
func ActivityOrganizerList(c *gin.Context) {
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
		response.ResponseErrorWithMsg(c, response.ParamFail, "query解析失败")
		return
	}
	cr.Label = "ActivityRulerList"
	var ResAllMsgList = make([]JoinAudit.ResList, 0)
	ResAllMsgList, err = JoinAudit.ResListWithJSON(cr)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, err.Error())
		return
	}
	if !cr.All {
		_, msg, _ := mysql.OpenActivityStates()
		response.ResponseSuccessWithMsg(c, msg, ResAllMsgList)
		return
	}
	response.ResponseSuccess(c, ResAllMsgList)
}

// 组织部材料审核
func ActivityOrganizerMaterialManager(c *gin.Context) {
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
		response.ResponseSuccessWithMsg(c, err.Error(), []struct{}{})
		return
	}
	resList := JoinAudit.IsPassWithJSON(cr, "organizer_material_is_pass")
	response.ResponseSuccess(c, resList)
}

// 组织部获取成绩列表
func ActivityOrganizerTrainList(c *gin.Context) {
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
		response.ResponseErrorWithMsg(c, response.ParamFail, "json数据解析失败")
		return
	}
	cr.Label = "ActivityOrganizerTrainList"
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

// ActivityOrganizerTrainManager 组织部考核成绩审核
func ActivityOrganizerTrainManager(c *gin.Context) {
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
		response.ResponseErrorWithMsg(c, response.ParamFail, "query数据解析失败")
		return
	}
	resList := JoinAudit.IsPassWithJSON(cr, "organizer_train_is_pass")
	response.ResponseSuccess(c, resList)
}

// 组织部更新分数
func SaveTrainScore(c *gin.Context) {
	token := token2.NewToken(c)
	_, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
		return
	}

	var cr []JoinAudit.Score
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "json数据解析失败")
		return
	}
	resList, err := JoinAudit.UpdateScore(cr)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, err.Error())
		return
	}
	response.ResponseSuccess(c, resList)
}
