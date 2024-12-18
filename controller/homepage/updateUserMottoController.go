package homepage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/pkg/response"
	"studentGrow/service"
	token2 "studentGrow/utils/token"
)

func UpdateHomepageMottoControl(c *gin.Context) {
	// 接收
	input := struct {
		UserMotto string `json:"user_motto" binding:"required"`
	}{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error(err.Error())
		return
	}

	// 接收
	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error("token错误")
	}

	// 业务
	err = service.UpdateHomepageMottoService(user.Username, input.UserMotto)
	if err != nil {
		response.ResponseError(c, response.ServerErrorCode)
		zap.L().Error(err.Error())
		return
	}

	// 响应
	response.ResponseSuccess(c, struct{}{})

}
