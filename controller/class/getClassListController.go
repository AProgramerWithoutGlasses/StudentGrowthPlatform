package class

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/models/jrx_model"
	"studentGrow/pkg/response"
	"studentGrow/service"
)

// 拿到所有班级的列表
func GetClassListControl(c *gin.Context) {
	//token := c.GetHeader("token")
	//username, err := token2.GetUsername(token)
	//if err != nil {
	//	response.ResponseError(c, response.ParamFail)
	//	zap.L().Error(err.Error())
	//	return
	//}

	// 业务
	classList, err := service.GetClassListService()
	if err != nil {
		response.ResponseError(c, response.ServerErrorCode)
		zap.L().Error(err.Error())
		return
	}

	output := struct {
		ClassList []jrx_model.Class `json:"class_list"`
	}{
		ClassList: classList,
	}
	response.ResponseSuccess(c, output)
}
