package class

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/pkg/response"
	"studentGrow/service"
)

func GetClassByGradeControl(c *gin.Context) {
	//token := c.GetHeader("token")
	//username, err := token2.GetUsername(token)
	//if err != nil {
	//	response.ResponseError(c, response.ParamFail)
	//	zap.L().Error(err.Error())
	//	return
	//}

	// 接收
	input := struct {
		Grade int `json:"grade"`
	}{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error(err.Error())
		return
	}

	// 业务
	classList, err := service.GetClassByGradeService(input.Grade)
	if err != nil {
		response.ResponseError(c, response.ServerErrorCode)
		zap.L().Error(err.Error())
		return
	}

	output := struct {
		GradeList []string `json:"grade_list"`
	}{
		GradeList: classList,
	}
	response.ResponseSuccess(c, output)
}
