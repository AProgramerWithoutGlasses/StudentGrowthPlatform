package routesJoinAudit

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	"studentGrow/models/gorm_model"
	"studentGrow/pkg/response"
	"studentGrow/service/JoinAudit"
	token2 "studentGrow/utils/token"
)

type StuFileMsg struct {
	ID       uint
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	Note     string `json:"note"`
}
type StuFileDelMsg struct {
	ID        int
	IsSuccess bool `json:"is_success"`
}

func GetStuFile(c *gin.Context) {
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
		return
	}
	var resMsg StuFileMsg
	var resList []StuFileMsg
	ActivityIsOpen, Msg, ActivityMsg := mysql.OpenActivityMsg()
	if !ActivityIsOpen {
		response.ResponseErrorWithMsg(c, response.ParamFail, Msg)
		return
	}

	// 判断上传或更新材料的前置条件是否符合
	stuFromMsg, _ := mysql.GetStuFromMsg(user.Username, ActivityMsg.ID)
	if stuFromMsg.ClassIsPass != "true" || stuFromMsg.RulerIsPass == "false" || stuFromMsg.OrganizerMaterialIsPass == "true" {
		response.ResponseErrorWithMsg(c, response.ParamFail, "前置条件不符合或者审核已通过")
		return
	}

	var fileList []gorm_model.JoinAuditFile
	mysql.DB.Where("username = ? AND join_audit_duty_id = ?", user.Username, ActivityMsg.ID).Find(&fileList)
	if len(fileList) == 0 {
		response.ResponseSuccessWithMsg(c, "用户文件不存在", resList)
		return
	}
	for _, file := range fileList {
		resMsg.ID = file.ID
		resMsg.FileName = file.FileName
		resMsg.FilePath = file.FilePath
		resMsg.Note = file.Note
		resList = append(resList, resMsg)
	}
	response.ResponseSuccess(c, resList)
	return
}

// SaveStuFile 保存提交文件
func SaveStuFile(c *gin.Context) {
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
		return
	}
	ActivityIsOpen, Msg, ActivityMsg := mysql.OpenActivityMsg()
	if !ActivityIsOpen {
		response.ResponseErrorWithMsg(c, response.ParamFail, Msg)
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ServerErrorCode, err.Error())
	}
	var resList []JoinAudit.StuFileUpload
	if form.File["material"] == nil || form.File["application"] == nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "无文件需要处理")
		return
	}
	//删除之前存在的文件

	err = mysql.DelUserFile(user.Username, ActivityMsg.ID)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "material 旧文件删除失败")
		return
	}

	//获取传入的文件
	fileList, ok := form.File["material"]
	if ok {
		for _, file := range fileList {
			resMsg := JoinAudit.FileUpload(file, user, ActivityMsg, "material")
			resList = append(resList, resMsg)
		}
	}
	fileList, ok = form.File["application"]
	if ok {
		for _, file := range fileList {
			resMsg := JoinAudit.FileUpload(file, user, ActivityMsg, "application")
			resList = append(resList, resMsg)
		}
	}
	if len(resList) == 0 {
		response.ResponseErrorWithMsg(c, response.ParamFail, "文件获取失败")
		return
	}
	mysql.DB.Model(&gorm_model.JoinAudit{}).Where("username = ? AND join_audit_duty_id = ?", user.Username, ActivityMsg.ID).Update("organizer_material_is_pass", "null")
	response.ResponseSuccess(c, resList)
}

// DelStuFile 文件删除
func DelStuFile(c *gin.Context) {
	type DelFileList struct {
		ID []int
	}
	token := token2.NewToken(c)
	_, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
		return
	}
	var cr DelFileList
	if err := c.ShouldBindQuery(&cr); err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "客户端数据解析失败")
		return
	}
	//判断删除时间是否合法
	ActivityIsOpen, Msg, _ := mysql.OpenActivityMsg()
	if !ActivityIsOpen {
		response.ResponseErrorWithMsg(c, response.ParamFail, Msg)
		return
	}
	var resFileDelMsgList []StuFileDelMsg
	for _, fileID := range cr.ID {
		var FileDelMsg StuFileDelMsg
		FileDelMsg.ID = fileID
		count := mysql.DelFileWithID(fileID)
		if count != 1 {
			FileDelMsg.IsSuccess = false
			resFileDelMsgList = append(resFileDelMsgList, FileDelMsg)
			continue
		}
		FileDelMsg.IsSuccess = true
		resFileDelMsgList = append(resFileDelMsgList, FileDelMsg)
	}
	response.ResponseSuccess(c, resFileDelMsgList)
}
