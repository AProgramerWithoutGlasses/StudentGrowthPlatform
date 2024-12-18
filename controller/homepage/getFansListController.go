package homepage

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"studentGrow/models/jrx_model"
	"studentGrow/pkg/response"
	"studentGrow/service"
	token2 "studentGrow/utils/token"
)

func GetFansListControl(c *gin.Context) {
	// 接收
	input := struct {
		Username string `json:"username" binding:"required"`
	}{}
	err := c.BindJSON(&input)
	if err != nil {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error(err.Error())
		return
	}

	///////
	//token := c.GetHeader("token")
	//tokenUsername, err := token2.GetUsername(token)
	//if err != nil {
	//	response.ResponseError(c, response.ParamFail)
	//	zap.L().Error(err.Error())
	//	return
	//}

	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error("token错误")
	}

	// 业务
	userfans, err := service.GetFansListService(input.Username, user.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

		} else {
			response.ResponseError(c, response.ServerErrorCode)
			zap.L().Error(err.Error())
			return
		}
	}

	// 响应
	output := struct {
		Userfans []jrx_model.HomepageFanStruct `json:"user_fans"`
	}{
		Userfans: userfans,
	}

	response.ResponseSuccess(c, output)

}
