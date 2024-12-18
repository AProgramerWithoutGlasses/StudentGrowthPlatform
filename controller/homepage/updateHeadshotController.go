package homepage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/pkg/response"
	"studentGrow/service"
	token2 "studentGrow/utils/token"
)

func UpdateHeadshotControl(c *gin.Context) {
	//input := struct {
	//	UserHeadshot string `form:"user_headshot"`
	//}{}
	//err := c.Bind(&input)
	//if err != nil {
	//	response.ResponseError(c, response.ParamFail)
	//	zap.L().Error(err.Error())
	//	return
	//}

	// 获取上传的文件
	file, err := c.FormFile("file")
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
	url, err := service.UpdateHeadshotService(file, user.Username)
	if err != nil {
		response.ResponseError(c, response.ServerErrorCode)
		zap.L().Error(err.Error())
		return
	}

	output := struct {
		UserHeadshot string `json:"user_headshot"`
	}{
		UserHeadshot: url,
	}

	// 响应
	response.ResponseSuccess(c, output)
}
