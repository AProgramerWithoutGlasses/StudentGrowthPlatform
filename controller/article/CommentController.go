package article

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	myErr "studentGrow/pkg/error"
	res "studentGrow/pkg/response"
	"studentGrow/service/comment"
	utils "studentGrow/utils/readMessage"
	"studentGrow/utils/timeConverter"
	"studentGrow/utils/token"
)

// PostCom 发布评论
func PostCom(c *gin.Context) {

	in := struct {
		TarUsername    string `json:"tar_username"`
		CommentType    int    `json:"comment_type"`
		CommentContent string `json:"comment_content"`
		Id             int    `json:"id"`
	}{}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		zap.L().Error("PostCom() controller.article.ShouldBindJSON err=", zap.Error(err))
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

	//新增评论
	err = comment.PostComment(username, in.TarUsername, in.CommentContent, in.Id, in.CommentType)

	if err != nil {
		zap.L().Error("PostCom() controller.article.PostComment err=", zap.Error(err))
		return
	}
	res.ResponseSuccess(c, struct{}{})
}

// GetLel1CommentsController 获取一级评论
func GetLel1CommentsController(c *gin.Context) {
	var input struct {
		Aid      int    `json:"article_id"`
		SortWay  string `json:"comment_sort"`
		Limit    int    `json:"comment_count"`
		Page     int    `json:"comment_page"`
		Username string `json:"username"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		zap.L().Error("GetLel1CommentsController() controller.article.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	comments, err := comment.GetLel1CommentsService(input.Aid, input.Limit, input.Page, input.Username, input.SortWay)
	if err != nil {
		zap.L().Error("GetLel1CommentsController() controller.article.GetLel1CommentsService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	for i := 0; i < len(comments); i++ {
		fmt.Println(i, comments[i].LikeAmount)
	}

	// 获取文章评论数
	commentNum, err := mysql.QueryArticleCommentNum(input.Aid)
	if err != nil {
		zap.L().Error("GetLel1CommentsController() controller.article.QueryArticleCommentNum err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	list := make([]map[string]any, 0)
	for _, comt := range comments {
		num, err := mysql.QuerySonCommentNum(int(comt.ID))
		if err != nil {
			zap.L().Error("GetLel1CommentsController() controller.article.QuerySonCommentNum err=", zap.Error(err))
			myErr.CheckErrors(err, c)
			return
		}
		list = append(list, map[string]any{
			"user_headshot":    comt.User.HeadShot,
			"username":         comt.User.Username,
			"name":             comt.User.Name,
			"comment_time":     timeConverter.IntervalConversion(comt.CreatedAt),
			"comment_content":  comt.Content,
			"id":               comt.ID,
			"comment_like_num": comt.LikeAmount,
			"comment_son_num":  num,
			"comment_if_like":  comt.IsLike,
			"p_id":             comt.Pid,
		})
	}

	res.ResponseSuccess(c, map[string]any{
		"comment_list": list,
		"comment_num":  commentNum,
	})
}

// GetSonCommentsController 获取子评论列表
func GetSonCommentsController(c *gin.Context) {
	var input struct {
		Cid      int    `json:"comment_id"`
		Username string `json:"username"`
		Limit    int    `json:"limit"`
		Page     int    `json:"page"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		zap.L().Error("GetSonCommentsController() controller.article.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	comments, err := comment.GetLelSonCommentListService(input.Cid, input.Limit, input.Page, input.Username)

	if err != nil {
		zap.L().Error("GetSonCommentsController() controller.article.GetLelSonCommentListService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	list := make([]map[string]any, 0)
	for _, comt := range comments {
		list = append(list, map[string]any{
			"user_headshot":    comt.User.HeadShot,
			"username":         comt.User.Username,
			"name":             comt.User.Name,
			"comment_time":     timeConverter.IntervalConversion(comt.CreatedAt),
			"comment_content":  comt.Content,
			"id":               comt.ID,
			"comment_like_num": comt.LikeAmount,
			"comment_if_like":  comt.IsLike,
			"p_id":             comt.Pid,
		})
	}

	res.ResponseSuccess(c, map[string]any{
		"comment_se_list": list,
	})
}

// DeleteCommentController 删除评论
func DeleteCommentController(c *gin.Context) {
	// 读取前端数据
	json, err := utils.GetJsonvalue(c)
	if err != nil {
		zap.L().Error("DeleteCommentController() controller.article.GetJsonvalue err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	// 获取评论id
	cid, err := json.GetInt("comment_id")
	if err != nil {
		zap.L().Error("DeleteCommentController() controller.article.GetInt err=", zap.Error(err))
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

	if err != nil {
		zap.L().Error("DeleteCommentController() controller.article.GetRole err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}
	// 删除评论
	err = comment.DeleteCommentService(cid, username)
	if err != nil {
		zap.L().Error("DeleteCommentController() controller.article.DeleteComment err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}
	res.ResponseSuccess(c, struct{}{})
}
