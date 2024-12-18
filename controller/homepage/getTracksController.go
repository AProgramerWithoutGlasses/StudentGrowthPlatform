package homepage

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/models/jrx_model"
	"studentGrow/pkg/response"
	"studentGrow/service"
	token2 "studentGrow/utils/token"
)

func GetTracksControl(c *gin.Context) {
	// 接收
	input := struct {
		Page  int `json:"page" binding:"required"`
		Limit int `json:"limit" binding:"required"`
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
	Tracks, err := service.GetTracksService(input.Page, input.Limit, user.Username)
	if err != nil {
		response.ResponseError(c, response.ServerErrorCode)
		zap.L().Error(err.Error())
		return
	}

	// 响应
	output := struct {
		InterInfo []jrx_model.HomepageTrack `json:"inter_info"`
	}{
		InterInfo: Tracks,
	}

	response.ResponseSuccess(c, output)
}
