package routesJoinAudit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	"studentGrow/models/gorm_model"
	"studentGrow/pkg/response"
	token2 "studentGrow/utils/token"
	"time"
)

type ResDelMsg struct {
	ID        int    `json:"ID"`
	IsSuccess bool   `json:"is_success"`
	Msg       string `json:"msg"`
}
type DelActivityIDList struct {
	ID []int `from:"ID"`
}
type ListResponse struct {
	List  any   `json:"list"`
	Count int64 `json:"count"`
}
type ReceiveActivityMsg struct {
	ID                    uint   `json:"ID"`
	ActivityName          string `json:"activity_name"`
	StartTime             string `json:"start_time"`
	StopTime              string `json:"stop_time"`
	PersonInCharge        string `json:"person_in_charge"`
	RulerName             string `json:"ruler_name"`
	OrganizerMaterialName string `json:"organizer_material_name"`
	OrganizerTrainName    string `json:"organizer_train_name"`
	IsShow                string `json:"is_show"`
	Note                  string `json:"note"`
}

func GetActivityList(c *gin.Context) {
	token := token2.NewToken(c)
	_, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
		return
	}
	var pagMsg mysql.Pagination
	err := c.ShouldBindJSON(&pagMsg)
	fmt.Println(pagMsg)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "json解析失败")
		return
	}
	pagMsg.Label = "ActivityList"
	ActivityMsgList, count, err := mysql.ComList(gorm_model.JoinAuditDuty{}, pagMsg)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "查询列表出现错误")
		return
	}
	response.ResponseSuccess(c, ListResponse{
		ActivityMsgList,
		count,
	})
}

// SaveActivityMsg 保存和更新活动
func SaveActivityMsg(c *gin.Context) {
	token := token2.NewToken(c)
	_, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
		return
	}
	var cr ReceiveActivityMsg
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "json解析失败")
		return
	}
	//转化时间
	startTime, err := time.Parse("2006-01-02T15:04:05.000-07:00", cr.StartTime)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "开始时间解析失败")
		return
	}
	stopTime, err := time.Parse("2006-01-02T15:04:05.000-07:00", cr.StopTime)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "结束时间解析失败")
		return
	}
	if cr.StartTime != "" && cr.StopTime != "" {
		if stopTime.Before(startTime) {
			response.ResponseErrorWithMsg(c, response.ParamFail, "结束时间早于开始时间")
			return
		}
	}

	var activityMsg = gorm_model.JoinAuditDuty{
		ActivityName:          cr.ActivityName,
		StartTime:             startTime,
		StopTime:              stopTime,
		PersonInCharge:        cr.PersonInCharge,
		RulerName:             cr.RulerName,
		OrganizerMaterialName: cr.OrganizerMaterialName,
		OrganizerTrainName:    cr.OrganizerTrainName,
		IsShow:                "false",
		Note:                  cr.Note,
	}

	//ID为0时创建新的活动
	if cr.ID == 0 {
		err = mysql.CreatActivity(activityMsg)
		if err != nil {
			fmt.Println(err.Error())
			response.ResponseErrorWithMsg(c, response.ParamFail, "活动创建失败")
			return
		} else {
			response.ResponseSuccessWithMsg(c, "活动创建成功", struct{}{})
			return
		}
	}
	//根据查询到的活动数量判断活动是否存在
	count := mysql.GetActivityNumberWithID(int(cr.ID))
	if count != 1 {
		response.ResponseErrorWithMsg(c, response.ParamFail, "查询活动异常")
		return
	}
	//存在记录时更新数据
	activityMsg.IsShow = cr.IsShow
	if cr.IsShow == "true" {
		err = mysql.CloseAllActivity()
		if err != nil {
			response.ResponseErrorWithMsg(c, response.ParamFail, "关闭其他活动异常")
			return
		}
	}
	err = mysql.UpdateActivityWithID(activityMsg, cr.ID)
	if err != nil {
		fmt.Println(err.Error())
		response.ResponseErrorWithMsg(c, response.ParamFail, "活动更新失败")
		return
	}
	response.ResponseSuccessWithMsg(c, "活动更新成功", struct{}{})
}

func DelActivityMsg(c *gin.Context) {
	token := token2.NewToken(c)
	_, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
		return
	}
	var cr DelActivityIDList
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "json数据解析失败")
		return
	}
	if len(cr.ID) == 0 {
		response.ResponseErrorWithMsg(c, response.ParamFail, "没有需要处理的数据")
		return
	}
	var resList []ResDelMsg
	var resDelMsg ResDelMsg
	for _, ActivityID := range cr.ID {
		resDelMsg.IsSuccess = false
		resDelMsg.ID = ActivityID
		//判断活动是否存在
		count := mysql.GetActivityNumberWithID(ActivityID)
		if count != 1 {
			resDelMsg.Msg = "活动查询异常"
			resList = append(resList, resDelMsg)
			continue
		}
		//删除活动
		err = mysql.DelActivityWithID(ActivityID)
		if err != nil {
			resDelMsg.Msg = "活动删除失败"
			resList = append(resList, resDelMsg)
			continue
		}
		resDelMsg.Msg = "活动删除成功"
		resDelMsg.IsSuccess = true
		resList = append(resList, resDelMsg)
	}
	response.ResponseSuccess(c, resList)
	return
}
