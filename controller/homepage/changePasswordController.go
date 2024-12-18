package homepage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/pkg/response"
	"studentGrow/service"
	token2 "studentGrow/utils/token"
)

func ChangePasswordControl(c *gin.Context) {
	// 接收
	input := struct {
		OldPwd string `json:"old_pwd" binding:"required"`
		NewPwd string `json:"new_pwd" binding:"required"`
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
	err = service.ChangePasswordService(user.Username, input.OldPwd, input.NewPwd)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ServerErrorCode, err.Error())
		zap.L().Error(err.Error())
		return
	}

	// 响应
	response.ResponseSuccess(c, struct {
	}{})

}
