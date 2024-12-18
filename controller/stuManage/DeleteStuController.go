package stuManage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/models/jrx_model"
	"studentGrow/pkg/response"
	"studentGrow/service"
	token2 "studentGrow/utils/token"
)

// 删除选中用户
func DeleteStuControl(c *gin.Context) {
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
	}

	// 接收请求
	var input struct {
		Selected_students []jrx_model.StuMesStruct `json:"selected_students"`
	}
	err := c.Bind(&input)
	if err != nil {
		response.ResponseErrorWithMsg(c, 500, "stuManager.DeleteStuControl() c.Bind() failed : "+err.Error())
		zap.L().Error("stuManager.DeleteStuControl() c.Bind() failed : ", zap.Error(err))
		return
	}

	// 业务
	deletedStuName, err := service.DeleteStuService(user.Username, input.Selected_students)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ServerErrorCode, err.Error())
		zap.L().Error(err.Error())
		return
	}

	// 响应成功
	response.ResponseSuccessWithMsg(c, deletedStuName+"删除成功!", nil)
}
