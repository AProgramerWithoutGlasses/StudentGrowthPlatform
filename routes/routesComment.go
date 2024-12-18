package routes

import (
	"github.com/gin-gonic/gin"
	"studentGrow/controller/article"
	"studentGrow/utils/middleWare"
	"studentGrow/utils/token"
)

func routesComment(r *gin.Engine) {
	ct := r.Group("/comment")
	// 获取一级评论
	ct.POST("/get_lel1comment", article.GetLel1CommentsController)
	// 获取子评论
	ct.POST("/get_lel2comment", article.GetSonCommentsController)
	// 删除评论
	ct.POST("/delete", middleWare.CORSMiddleware(), token.AuthMiddleware(), article.DeleteCommentController)
}
