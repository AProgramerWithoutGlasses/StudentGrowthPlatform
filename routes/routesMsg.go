package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/controller/message"
	"studentGrow/dao/mysql"
	"studentGrow/models/casbinModels"
	"studentGrow/utils/middleWare"
	"studentGrow/utils/token"
)

func routesMsg(r *gin.Engine) {
	gp := r.Group("/report_box")

	casbinService, err := casbinModels.NewCasbinService(mysql.DB)
	if err != nil {
		zap.L().Error("routesMsg() routes.routesArticle.NewCasbinService err=", zap.Error(err))
		return
	}

	// 查看举报信息
	gp.POST("/getlist", token.AuthMiddleware(), middleWare.NewCasbinAuth(casbinService), message.GetUnreadReportsController)

	// 确认举报信息
	gp.POST("/ack", token.AuthMiddleware(), message.AckUnreadReportsController)

	msg := r.Group("/message")

	// 获取系统消息
	msg.POST("/get_system", token.AuthMiddleware(), message.GetSystemMsgController)

	// 获取管理员消息
	msg.POST("/get_manager", token.AuthMiddleware(), message.GetManagerMsgController)

	// 获取点赞消息
	msg.POST("/get_thumbList", token.AuthMiddleware(), message.GetLikeMsgController)

	// 获取收藏消息
	msg.POST("/get_starList", token.AuthMiddleware(), message.GetCollectMsgController)

	// 获取评论消息
	msg.POST("/get_comList", token.AuthMiddleware(), message.GetCommentMsgController)

	// 确认互动消息
	msg.POST("/ack_interactMsg", token.AuthMiddleware(), message.AckInterMsgController)

	// 确认系统消息
	msg.POST("/ack_systemMsg", token.AuthMiddleware(), message.AckSystemMsgController)

	// 确认管理员消息
	msg.POST("/ack_managerMsg", token.AuthMiddleware(), message.AckManagerMsgController)

	// 发布管理员通知
	msg.POST("/publish_managerMsg", token.AuthMiddleware(), message.PublishManagerMsgController)

	// 发布系统通知
	msg.POST("/publish_systemMsg", token.AuthMiddleware(), message.PublishSystemMsgController)

	// 删除系统通知
	msg.POST("/delete_systemMsg", token.AuthMiddleware(), message.DeleteSystemMsgController)

	// 删除管理员通知
	msg.POST("/delete_managerMsg", token.AuthMiddleware(), message.DeleteManagerMsgController)

	// 一键已读互动消息
	msg.POST("/ack_interactAllMsg", token.AuthMiddleware(), message.AckAllInterMsgController)
}
