package stuManage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/pkg/response"
	"studentGrow/service"
)

func ReSetPasswordControl(c *gin.Context) {
	// 接收
	input := struct {
		Username string `json:"username" binding:"required"`
	}{}
	err := c.Bind(&input)
	if err != nil {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error(err.Error())
		return
	}

	// 业务
	err = service.ReSetPasswordService(input.Username)
	if err != nil {
		zap.L().Error(err.Error())
		response.ResponseError(c, response.ServerErrorCode)
		return
	}

	// 输出
	response.ResponseSuccessWithMsg(c, "密码已成功重置为：123456", struct{}{})

}
