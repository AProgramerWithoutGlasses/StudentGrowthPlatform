package teacherManage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/models/gorm_model"
	"studentGrow/pkg/response"
	"studentGrow/service"
)

func DeleteTeacherControl(c *gin.Context) {
	// 接收请求
	var input gorm_model.User
	err := c.Bind(&input)
	if err != nil {
		response.ResponseErrorWithMsg(c, 500, "teacherManager.DeleteStuControl() c.Bind() failed : "+err.Error())
		zap.L().Error("teacherManager.DeleteStuControl() c.Bind() failed : ", zap.Error(err))
		return
	}

	// 业务
	err = service.DeleteSingleTeacherService(input)
	if err != nil {
		response.ResponseErrorWithMsg(c, 500, "teacherManager.DeleteStuControl() service.DeleteSingleTeacherService() failed : "+err.Error())
		zap.L().Error("teacherManager.DeleteStuControl() service.DeleteSingleTeacherService() failed : ", zap.Error(err))
		return
	}

	// 响应成功
	response.ResponseSuccess(c, "删除成功!")
}
