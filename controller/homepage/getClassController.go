package homepage

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"studentGrow/models/jrx_model"
	"studentGrow/pkg/response"
	"studentGrow/service"
)

func GetClassControl(c *gin.Context) {
	// 接收
	//token := c.GetHeader("token")
	//username, err := token2.GetUsername(token)
	//if err != nil {
	//	response.ResponseError(c, response.ParamFail)
	//	zap.L().Error(err.Error())
	//	return
	//}

	input := struct {
		Username string `json:"username" binding:"required"`
	}{}
	err := c.BindJSON(&input)
	if err != nil {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error(err.Error())
		return
	}

	// 业务
	classmateList, err := service.GetClassmateListService(input.Username)
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
		History []jrx_model.HomepageClassmateStruct `json:"student_total"`
	}{
		History: classmateList,
	}

	response.ResponseSuccess(c, output)
}
