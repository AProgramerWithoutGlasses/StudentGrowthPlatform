package stuManage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/models/jrx_model"
	"studentGrow/pkg/response"
	"studentGrow/service"
	token2 "studentGrow/utils/token"
)

type InnerInput struct {
	Username string `json:"username" binding:"required"`
	Year     string `json:"year" binding:"required"`
}

type Input struct {
	Student     InnerInput `json:"student"`
	ManagerType string     `json:"managerType"`
}

// 设置用户为管理员
func SetStuManagerControl(c *gin.Context) {
	// 接收
	var setStuManagerModel jrx_model.SetStuManagerModel
	err := c.Bind(&setStuManagerModel)
	if err != nil {
		zap.L().Error("stuManage.SetStuManagerControl() c.Bind() failed", zap.Error(err))
		response.ResponseError(c, response.ParamFail)
		return
	}

	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
	}

	// 业务
	err = service.SetStuManagerService(setStuManagerModel.Student.Username, user.Username, setStuManagerModel.ManagerType, setStuManagerModel.Student.Year)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ServerErrorCode, err.Error())
		zap.L().Error("stuManager.SetStuManagerControl() service.SetStuManagerService() failed : ", zap.Error(err))
		return
	}

	// 响应
	response.ResponseSuccessWithMsg(c, "已将用户 "+setStuManagerModel.Student.Username+" 设置为 "+setStuManagerModel.ManagerType, "")

}
