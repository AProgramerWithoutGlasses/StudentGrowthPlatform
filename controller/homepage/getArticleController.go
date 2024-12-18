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

func GetArticleControl(c *gin.Context) {
	// 接收
	input := struct {
		Page     int    `form:"page" binding:"required"`
		Limit    int    `form:"limit" binding:"required"`
		Username string `form:"username" binding:"required"`
	}{}
	err := c.ShouldBindQuery(&input)
	if err != nil {
		response.ResponseError(c, response.ParamFail)
		zap.L().Error(err.Error())
		return
	}

	token := token2.NewToken(c)
	user, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
		return
	}

	// 业务
	articleList, err := service.GetArticleService(input.Page, input.Limit, input.Username, user.Username)
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
		Content []jrx_model.HomepageArticleHistoryStruct `json:"content"`
	}{
		Content: articleList,
	}

	response.ResponseSuccess(c, output)

}
