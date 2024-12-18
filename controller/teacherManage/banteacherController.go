package teacherManage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/pkg/response"
	"studentGrow/service"
)

func BanTeacherControl(c *gin.Context) {
	// 接收参数
	input := struct {
		Username string `json:"username" binding:"required"`
	}{}
	err := c.Bind(&input)
	if err != nil {
		response.ResponseErrorWithMsg(c, 500, "teacherManager.BanTeacherControl() c.Bind() failed : "+err.Error())
		zap.L().Error("teacherManager.BanTeacherControl() c.Bind() failed : ", zap.Error(err))
		return
	}

	// 业务
	name, temp, err := service.BanUserService(input.Username)
	if err != nil {
		response.ResponseErrorWithMsg(c, 500, "teacherManager.BanTeacherControl() service.BanUserService() failed : "+err.Error())
		zap.L().Error("teacherManager.BanTeacherControl() service.BanUserService() failed : ", zap.Error(err))
		return
	}

	// 响应
	if temp == 1 {
		response.ResponseSuccess(c, "已将用户"+name+"封禁")
	} else if temp == 0 {
		response.ResponseSuccess(c, "已将用户"+name+"取消封禁")
	}

}
