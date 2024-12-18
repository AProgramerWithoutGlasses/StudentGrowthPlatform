package teacherManage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	"studentGrow/pkg/response"
)

func GetUsername(c *gin.Context) {
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
	id, err := mysql.GetIdByUsername(input.Username)
	if err != nil {
		zap.L().Error(err.Error())
		response.ResponseError(c, response.ServerErrorCode)
		return
	}

	// 输出
	output := struct {
		Id int `json:"id"`
	}{
		Id: id,
	}
	response.ResponseSuccess(c, output)

}
