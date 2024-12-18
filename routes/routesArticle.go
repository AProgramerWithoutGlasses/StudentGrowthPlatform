package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/controller/article"
	"studentGrow/dao/mysql"
	"studentGrow/models/casbinModels"
	"studentGrow/utils/middleWare"
	"studentGrow/utils/token"
)

func routesArticle(r *gin.Engine) {
	casbinService, err := casbinModels.NewCasbinService(mysql.DB)
	if err != nil {
		zap.L().Error("routesArticle() routes.routesArticle.NewCasbinService err=", zap.Error(err))
		return
	}
	at := r.Group("/article")
	// 获取文章内容
	at.POST("/content", article.GetArticleIdController)
	// 获取文章列表
	at.POST("/list", token.AuthMiddleware(), article.GetArticleListController)
	// 对文章进行评论
	at.POST("/comment", token.AuthMiddleware(), article.PostCom)
	// 获取文章标签
	at.POST("/publish/get_tags", article.SendTopicTagsController)
	//文章或评论点赞
	at.POST("/like", token.AuthMiddleware(), article.LikeController)
	//封禁文章
	at.POST("/ban", token.AuthMiddleware(), middleWare.NewCasbinAuth(casbinService), article.BannedArticleController)
	//删除文章
	at.POST("/delete", token.AuthMiddleware(), article.DeleteArticleController)
	//举报文章
	at.POST("/report", token.AuthMiddleware(), article.ReportArticle)
	// 获取今日热帖
	at.POST("/hotpost/title", article.GetHotArticlesOfDayController)
	// 首页模糊搜索
	at.POST("/search_first", article.SelectArticleAndUserListByPageFirstPageController)
	// 收藏
	at.POST("/collect", token.AuthMiddleware(), article.CollectArticleController)
	// 发布文章
	at.POST("/publish", token.AuthMiddleware(), article.PublishArticleController)
	// 班级分类文章列表
	at.POST("/class_search", article.GetArticleByClassController)
	// 修改文章私密状态
	at.POST("/status", token.AuthMiddleware(), article.ReviseArticleStatusController)
	// 评选优秀帖子
	at.POST("/select_good_article", token.AuthMiddleware(), article.SelectGoodArticleController)
	// 帖子高级筛选
	at.POST("/filter", article.AdvancedArticleFilteringController)
	// 获取优秀帖子
	at.POST("/getGoodArticles", token.AuthMiddleware(), article.GetGoodArticlesController)
	// 取消收藏
	//at.POST("/cancel_collect", article.CancelCollectArticleController)
	// 查看收藏列表
	//at.POST("/get_collects", article.GetArticleListForSelectController)
	//获取文章点赞数量
	//at.POST("/like_nums", article.GetObjLikeNumController)
	//检查当前是否点赞
	//at.POST("/isLike", article.CheckLikeOrNotController)
}
