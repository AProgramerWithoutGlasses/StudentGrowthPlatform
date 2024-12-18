package teacherManage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/pkg/response"
	"studentGrow/service"
)

// 设置老师管理员
func SetTeacherManagerControl(c *gin.Context) {
	// 请求
	var input struct {
		Username    string `json:"username" binding:"required"`
		ManagerType string `json:"manager_type" binding:"required"`
	}
	err := c.Bind(&input)
	if err != nil {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error("teacherManage.SetTeacherManagerControl() c.Bind() failed : ", zap.Error(err))
		return
	}

	// 业务
	err = service.SetTeacherManagerService(input.Username, input.ManagerType)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ServerErrorCode, err.Error())
		zap.L().Error("teacherManage.SetTeacherManagerControl() service.SetTeacherManagerService() failed : ", zap.Error(err))
		return
	}

	// 响应
	response.ResponseSuccessWithMsg(c, "已将用户 "+input.Username+" 设置为 "+input.ManagerType, struct{}{})
}
