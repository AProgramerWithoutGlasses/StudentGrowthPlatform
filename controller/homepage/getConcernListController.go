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

func GetConcernListControl(c *gin.Context) {
	// 接收
	input := struct {
		Username string `json:"username" binding:"required"`
	}{}
	err := c.BindJSON(&input)
	if err != nil {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error(err.Error())
		return
	} //HH

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
		return
	}

	// 业务
	userConcern, err := service.GetConcernListService(input.Username, user.Username)
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
		UserConcern []jrx_model.HomepageFanStruct `json:"user_concern"`
	}{
		UserConcern: userConcern,
	}

	response.ResponseSuccess(c, output)

}
