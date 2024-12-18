package message

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	myErr "studentGrow/pkg/error"
	res "studentGrow/pkg/response"
	"studentGrow/service/message"
	readUtil "studentGrow/utils/readMessage"
	"studentGrow/utils/timeConverter"
	"studentGrow/utils/token"
)

// GetUnreadReportsController 获取未读举报信息
func GetUnreadReportsController(c *gin.Context) {

	aToken := token.NewToken(c)
	user, exist := aToken.GetUser()
	if !exist {
		res.ResponseError(c, res.TokenError)
		zap.L().Error("token错误")
		return
	}

	username := user.Username
	role, err := aToken.GetRole()
	if err != nil {
		zap.L().Error("GetUnreadReportsController() controller.message.GetRole err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	input := struct {
		Limit int `json:"limit"`
		Page  int `json:"page"`
	}{}

	err = c.ShouldBindJSON(&input)
	if err != nil {
		zap.L().Error("GetUnreadReportsController() controller.message.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	// 获取未读举报列表
	reports, err := message.GetUnreadReportsForService(username, role, input.Limit, input.Page)
	if err != nil {
		zap.L().Error("GetUnreadReportsController() controller.message.GetUnreadReportsForService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	// 获取未读举报数目
	count, err := message.GetUnreadReportNumForService(username, role)
	if err != nil {
		zap.L().Error("GetUnreadReportsController() controller.message.GetUnreadReportsForService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	// 通过文章id映射report_content,article_content
	var reportContent = make(map[uint][]map[string]any)
	var articleContent = make(map[uint]string)
	for _, item := range reports {
		reportContent[item.ArticleID] = append(reportContent[item.ArticleID], map[string]any{
			"report_time": timeConverter.IntervalConversion(item.CreatedAt),
			"report_msg":  item.Msg,
		})
		articleContent[item.ArticleID] = item.Article.Content
	}

	list := make([]map[string]any, 0)
	for key, val := range reportContent {
		list = append(list, map[string]any{
			"article_id":      key,
			"article_content": articleContent[key],
			"report_content":  val,
		})
	}

	// 返回响应
	res.ResponseSuccess(c, map[string]any{
		"article_ban":  list,
		"username":     username,
		"unread_count": count,
	})
}

// AckUnreadReportsController 确认未读举报信息
func AckUnreadReportsController(c *gin.Context) {
	// 通过token获取username
	aToken := token.NewToken(c)
	user, exist := aToken.GetUser()
	if !exist {
		res.ResponseError(c, res.TokenError)
		zap.L().Error("token错误")
		return
	}

	username := user.Username
	role, err := aToken.GetRole()
	if err != nil {
		zap.L().Error("AckUnreadReportsController() controller.message.GetRole err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	//获取前端发送的数据
	json, err := readUtil.GetJsonvalue(c)
	if err != nil {
		zap.L().Error("AckUnreadReportsController() controller.message.GetJsonvalue err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	reportId, err := json.GetInt("article_id")
	if err != nil {
		zap.L().Error("AckUnreadReportsController() controller.message.GetInt err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	// 确认未读举报消息
	err = message.AckUnreadReportsService(reportId, username, role)
	if err != nil {
		zap.L().Error("AckUnreadReportsController() controller.message.AckUnreadReportsService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	// 返回响应
	res.ResponseSuccess(c, struct{}{})
}
