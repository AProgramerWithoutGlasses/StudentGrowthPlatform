package article

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	myErr "studentGrow/pkg/error"
	res "studentGrow/pkg/response"
	"studentGrow/service/article"
	readUtil "studentGrow/utils/readMessage"
)

// AddTopicsController 添加话题
func AddTopicsController(c *gin.Context) {
	//获取前端发送的数据
	json, err := readUtil.GetJsonvalue(c)

	if err != nil {
		zap.L().Error("AddTopicsController() controller.article.getArticle.GetJsonvalue err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}
	err = article.AddTopicsService(json)
	if err != nil {
		zap.L().Error("AddTopicsController() controller.article.getArticle.AddTopicsService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	res.ResponseSuccess(c, struct{}{})

}

// GetAllTopicsController 获取所有话题
func GetAllTopicsController(c *gin.Context) {

	result, err := article.GetAllTopicsService()
	if err != nil {
		zap.L().Error("AddTopicsController() controller.article.getArticle.GetAllTopicsService err=", zap.Error(err))
		if err != nil {
			myErr.CheckErrors(err, c)
			return
		}
	}

	res.ResponseSuccess(c, map[string]any{
		"topic_list": result,
	})

}

// AddTagsByTopicController 添加标签
func AddTagsByTopicController(c *gin.Context) {
	in := struct {
		Topic string   `json:"topic"`
		Tags  []string `json:"tags"`
	}{}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		zap.L().Error("AddArticleTagsController() controller.article.getArticle.ShouldBindJSON err=", zap.Error(err))
		return
	}

	err = article.AddTagsByTopicService(in.Topic, in.Tags)
	if err != nil {
		zap.L().Error("AddArticleTagsController() controller.article.getArticle.AddTagsByTopicService err=", zap.Error(err))
		if err != nil {
			myErr.CheckErrors(err, c)
			return
		}
		res.ResponseError(c, res.ServerErrorCode)
		return
	}

	res.ResponseSuccess(c, struct{}{})
}

// GetTagsByTopicController 获取标签
func GetTagsByTopicController(c *gin.Context) {
	in := struct {
		TopicID int `json:"topic_id"`
	}{}
	err := c.ShouldBindJSON(&in)
	if err != nil {
		zap.L().Error("GetTagsByTopicController() controller.article.getArticle.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	result, err := article.GetTagsByTopicService(in.TopicID)
	if err != nil {
		zap.L().Error("GetTagsByTopicController() controller.article.getArticle.GetTagsByTopicService err=", zap.Error(err))
		if errors.Is(err, myErr.ErrNotFoundError) {
			myErr.CheckErrors(err, c)
			return
		}
		res.ResponseError(c, res.ServerErrorCode)
		return
	}
	res.ResponseSuccess(c, result)
}

// SendTopicTagsController 发送话题标签数据
func SendTopicTagsController(c *gin.Context) {
	in := struct {
		TopicID int `json:"topic_id"`
	}{}
	err := c.ShouldBindJSON(&in)
	if err != nil {
		zap.L().Error("SendTopicTagsController() controller.article.getArticle.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	//获取到查询的标签
	result, err := article.GetTagsByTopicService(in.TopicID)
	if err != nil {
		zap.L().Error("SendTopicTagsController() controller.article.getArticle.GetTagsByTopicService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	//返回响应
	res.ResponseSuccess(c, result)
}
