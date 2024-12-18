package article

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	myErr "studentGrow/pkg/error"
	res "studentGrow/pkg/response"
	"studentGrow/service/article"
	"studentGrow/utils/token"
)

// CollectArticleController 收藏文章
func CollectArticleController(c *gin.Context) {
	in := struct {
		ArticleId   int    `json:"article_id"`
		TarUsername string `json:"tar_username"`
	}{}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		zap.L().Error("CollectArticleController() controller.article.CollectController.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	// 通过token获取username
	aToken := token.NewToken(c)
	user, exist := aToken.GetUser()
	if !exist {
		res.ResponseError(c, res.TokenError)
		zap.L().Error("token错误")
		return
	}

	username := user.Username

	// 收藏
	err = article.CollectOrNotService(strconv.Itoa(in.ArticleId), username, in.TarUsername)
	if err != nil {
		zap.L().Error("CollectArticleController() controller.article.CollectController.CollectOrNotService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	res.ResponseSuccess(c, struct{}{})

}
