package routes

import (
	"github.com/gin-gonic/gin"
	"studentGrow/controller/article"
)

func routesTopic(r *gin.Engine) {
	topic := r.Group("/publish")
	// 添加话题
	topic.POST("/add_topic", article.AddTopicsController)
	// 获取话题
	topic.POST("/get_topic", article.GetAllTopicsController)
	// 添加标签
	topic.POST("/add_tags", article.AddTagsByTopicController)
	// 获取标签
	topic.POST("/get_tags", article.GetTagsByTopicController)
}
