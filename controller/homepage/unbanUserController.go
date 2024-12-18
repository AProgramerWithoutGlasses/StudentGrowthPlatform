package homepage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/pkg/response"
	"studentGrow/service"
	token2 "studentGrow/utils/token"
)

func UnbanUserControl(c *gin.Context) {
	// 接收数据
	input := struct {
		UnbanUsername string `json:"unban_username" binding:"required"`
	}{}
	err := c.BindJSON(&input)
	if err != nil {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error(err.Error())
		return
	}

	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error("token错误")
	}

	// 业务
	err = service.UnbanHomepageUserService(input.UnbanUsername, user.Username)
	if err != nil {
		response.ResponseError(c, response.ServerErrorCode)
		return
	}

	// 响应
	response.ResponseSuccess(c, struct{}{})
}
