package teacherManage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/models/jrx_model"
	"studentGrow/pkg/response"
	"studentGrow/service"
)

func EditTeacherControl(c *gin.Context) {
	// 接收数据
	var input jrx_model.ChangeTeacherMesStruct
	err := c.Bind(&input)
	if err != nil {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error("stuManage.EditStuControl() c.Bind() err : ", zap.Error(err))
		return
	}

	// 业务
	err = service.EditTeacherService(input)
	if err != nil {
		response.ResponseError(c, response.ServerErrorCode)
		zap.L().Error("stuManage.EditStuControl() service.EditTeacherService() err : ", zap.Error(err))
		return
	}

	// 响应
	response.ResponseSuccess(c, nil)
}
