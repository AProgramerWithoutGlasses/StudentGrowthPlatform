package homepage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/pkg/response"
	"studentGrow/service"
	token2 "studentGrow/utils/token"
)

func GetIsConcernControl(c *gin.Context) {
	input := struct {
		//Username      string `json:"username"`
		OtherUsername string `json:"other_username" binding:"required"`
	}{}

	err := c.BindJSON(&input)
	if err != nil {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error(err.Error())
		return
	}

	// 校验

	// 获取角色
	//token := c.GetHeader("token")
	//username, err := token2.GetUsername(token) // class, grade(1-4), collge, superman
	//if err != nil {
	//	response.ResponseError(c, response.ServerErrorCode)
	//	zap.L().Error(err.Error())
	//	return
	//}

	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error("token错误")
	}

	isConcern, err := service.GetIsConcernService(user.Username, input.OtherUsername)
	if err != nil {
		response.ResponseError(c, response.ServerErrorCode)
		zap.L().Error(err.Error())
		return
	}

	output := struct {
		IsConcern bool `json:"is_concern"`
	}{
		IsConcern: isConcern,
	}
	response.ResponseSuccess(c, output)

}
