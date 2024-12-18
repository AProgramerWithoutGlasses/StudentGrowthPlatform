package routesJoinAudit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	"studentGrow/pkg/response"
	token2 "studentGrow/utils/token"
)

// StuMsg 学生表单信息
type StuMsg struct {
	ActivityName            string   `json:"activity_name"`
	Username                string   `json:"username" `
	UserClass               string   `json:"user_class"`
	FromStuClass            string   `json:"from_stu_class"`
	Name                    string   `json:"name"`
	MoralCoin               float64  `json:"moral_coin"` //道德币
	ClassList               []string `json:"class_list"`
	ComprehensiveScore      float64  `json:"comprehensive_score"` //综测成绩
	ClassIsPass             string   `json:"class_is_pass"`
	RulerIsPass             string   `json:"ruler_is_pass" `              //纪权部综测成绩审核结果
	OrganizerMaterialIsPass string   `json:"organizer_material_is_pass" ` //组织部材料审核结果
}

// GetStudForm 获取用户的输入
func GetStudForm(c *gin.Context) {
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
	}

	ActivityIsOpen, Msg, ActivityMsg := mysql.OpenActivityMsg()
	if !ActivityIsOpen {
		response.ResponseErrorWithMsg(c, response.ParamFail, Msg)
		return
	}
	classList := mysql.AllClassList()
	var resStuMsg StuMsg
	//获取用户表单信息
	isExist, stuMsg := mysql.StuFormMsg(user.Username, ActivityMsg.ID)
	//初次进入未提交过表单状态
	if !isExist {
		resStuMsg = StuMsg{
			Username:                user.Username,
			Name:                    user.Name,
			UserClass:               "",         //需要填写的
			FromStuClass:            user.Class, //学生表中的
			ClassList:               classList,
			ActivityName:            ActivityMsg.ActivityName,
			ClassIsPass:             "null",
			RulerIsPass:             "null",
			OrganizerMaterialIsPass: "null",
		}
		response.ResponseSuccess(c, resStuMsg)
		return
	}

	resStuMsg = StuMsg{
		Username:                user.Username,
		UserClass:               stuMsg.UserClass,
		FromStuClass:            user.Class,
		ClassList:               classList,
		ActivityName:            ActivityMsg.ActivityName,
		Name:                    stuMsg.Name,
		MoralCoin:               stuMsg.MoralCoin,
		ComprehensiveScore:      stuMsg.ComprehensiveScore,
		ClassIsPass:             stuMsg.ClassIsPass,
		RulerIsPass:             stuMsg.RulerIsPass,
		OrganizerMaterialIsPass: stuMsg.OrganizerMaterialIsPass,
	}
	response.ResponseSuccess(c, resStuMsg)
}

// SaveStudForm 添加表单数据到数据库
func SaveStudForm(c *gin.Context) {
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	fmt.Println(user.Username)
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
	}
	var cr StuMsg
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "json解析失败")
		return
	}
	ActivityIsOpen, Msg, ActivityMsg := mysql.OpenActivityMsg()
	if !ActivityIsOpen {
		response.ResponseErrorWithMsg(c, response.ParamFail, Msg)
		return
	}

	isExist, stuMsg := mysql.StuFormMsg(user.Username, ActivityMsg.ID)
	stuMsg.Gender = user.Gender
	stuMsg.Name = cr.Name
	stuMsg.UserClass = cr.UserClass
	stuMsg.MoralCoin = cr.MoralCoin
	stuMsg.ComprehensiveScore = cr.ComprehensiveScore
	stuMsg.JoinAuditDuty = ActivityMsg
	stuMsg.ClassIsPass = "null"
	stuMsg.RulerIsPass = "null"
	stuMsg.OrganizerTrainIsPass = "null"
	stuMsg.OrganizerMaterialIsPass = "null"
	if isExist {
		err = mysql.UpdateStuForm(stuMsg)
		if err != nil {
			response.ResponseErrorWithMsg(c, response.ParamFail, "信息更新失败")
			return
		}
		response.ResponseSuccessWithMsg(c, "信息更新成功", struct{}{})
		return
	} else {
		stuMsg.Username = user.Username
		err = mysql.CreateStuForm(stuMsg)
		if err != nil {
			response.ResponseErrorWithMsg(c, response.ParamFail, "信息创建失败")
			return
		}
		response.ResponseSuccessWithMsg(c, "信息创建成功", struct{}{})
	}
}
