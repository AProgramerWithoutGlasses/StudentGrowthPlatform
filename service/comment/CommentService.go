package comment

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"studentGrow/dao/mysql"
	"studentGrow/dao/redis"
	"studentGrow/models/gorm_model"
	myErr "studentGrow/pkg/error"
	"studentGrow/pkg/sse"
	"studentGrow/service/article"
	NotificationPush "studentGrow/service/notificationPush"
	timeUtil "studentGrow/utils/timeConverter"
)

// PostComment 发布评论
func PostComment(username, tarUsername, content string, id, commentType int) error {
	//类型comment_type:‘article’or‘comment’;id;comment_content;comment_username

	//获取用户id
	uid, err := mysql.SelectUserByUsername(username)
	fmt.Println(uid)
	if err != nil {
		zap.L().Error("PostComment() service.article.SelectUserByUsername err=", zap.Error(err))
		return err
	}

	var cid int
	//判断评论类型

	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		switch commentType {
		//给文章评论
		case 0:
			//向数据库插入评论数据
			cid, err = mysql.InsertIntoCommentsForArticle(content, id, uid, tx)
			if err != nil {
				zap.L().Error("PostComment() service.article.InsertIntoCommentsForArticle err=", zap.Error(err))
				return err
			}

		case 1:
			//向数据库插入评论数据
			cid, err = mysql.InsertIntoCommentsForComment(content, uid, id, tx)
			if err != nil {
				zap.L().Error("PostComment() service.article.InsertIntoCommentsForComment err=", zap.Error(err))
				return err
			}
		default:
			return myErr.DataFormatError()
		}
		return nil
	})
	if err != nil {
		zap.L().Error("PostComment() service.article.Transaction err=", zap.Error(err))
		return err
	}

	// 获取评论的文章id
	aid, err := mysql.QueryArticleIdByCommentId(cid)
	if err != nil {
		zap.L().Error("PostComment() service.article.QueryArticleIdByCommentId err=", zap.Error(err))
		return err
	}
	// 增加文章评论数
	num, err := mysql.QueryArticleCommentNum(aid)
	if err != nil {
		zap.L().Error("PostComment() service.article.QueryArticleCommentNum err=", zap.Error(err))
		return err
	}
	err = mysql.UpdateArticleCommentNum(aid, num+1, mysql.DB)
	if err != nil {
		zap.L().Error("PostComment() service.article.UpdateArticleCommentNum err=", zap.Error(err))
		return err
	}

	// 将评论数据加入redis
	redis.RDB.HSet("comment", strconv.Itoa(cid), 0)

	notification, err := NotificationPush.BuildCommentNotification(username, tarUsername, id, commentType)
	if err != nil {
		zap.L().Error("PostComment() service.article.BuildCommentNotification err=", zap.Error(err))
		return err
	}
	sse.SendInterNotification(*notification)

	return nil
}

// GetLel1CommentsService 获取一级评论详情列表
func GetLel1CommentsService(aid, limit, page int, username, sortWay string) (gorm_model.Comments, error) {
	// 判断该用户是否被封禁或私密
	bl, err := article.QueryArticleStatusAndBanById(aid)
	if err != nil {
		zap.L().Error("GetLel1CommentsService() service.article.QueryArticleStatusAndBanById err=", zap.Error(err))
		return nil, err
	}
	if !bl {
		return nil, myErr.ErrNotFoundError
	}

	// 分页查询评论
	comments, err := mysql.QueryLevelOneComments(aid, limit, page)
	if err != nil {
		zap.L().Error("GetLel1CommentsService() service.article.QueryLevelOneComments err=", zap.Error(err))
		return nil, err
	}
	// 排序
	if sortWay == "hot" {
		sort.Sort(comments)
	}
	// 判断是否点赞, 并计算其子评论数量, 计算评论时间
	for i := 0; i < len(comments); i++ {
		liked, err := redis.IsUserLiked(strconv.Itoa(int(comments[i].ID)), username, 1)
		if err != nil {
			zap.L().Error("GetLel1CommentsService() service.article.IsUserLiked err=", zap.Error(err))
			return nil, err
		}
		comments[i].IsLike = liked

		num, err := mysql.QuerySonCommentNum(int(comments[i].ID))
		if err != nil {
			zap.L().Error("GetLel1CommentsService() service.article.QuerySonCommentNum err=", zap.Error(err))
			return nil, err
		}
		comments[i].ReplyCount = num

		comments[i].Time = timeUtil.IntervalConversion(comments[i].CreatedAt)
	}
	return comments, nil
}

// GetLelSonCommentListService 获取子评论列表
func GetLelSonCommentListService(cid, limit, page int, username string) ([]gorm_model.Comment, error) {
	// 获取文章对应的评论
	comments, err := mysql.QueryLevelSonComments(cid, limit, page)
	if err != nil {
		zap.L().Error("GetLelSonCommentListService() service.article.QueryLevelSonComments err=", zap.Error(err))
		return nil, err
	}

	// 该用户是否点赞, 计算评论时间
	for i := 0; i < len(comments); i++ {
		liked, err := redis.IsUserLiked(strconv.Itoa(int(comments[i].ID)), username, 1)
		if err != nil {
			zap.L().Error("GetLelSonCommentListService() service.article.IsUserLiked err=", zap.Error(err))
			return nil, err
		}
		comments[i].IsLike = liked

		comments[i].Time = timeUtil.IntervalConversion(comments[i].CreatedAt)
	}

	return comments, nil
}

// DeleteCommentService 删除评论
func DeleteCommentService(cid int, username string) error {
	comment, err := mysql.QueryCommentById(cid)
	if err != nil {
		zap.L().Error("DeleteCommentService() dao.mysql.nzx_sql.QueryCommentById err=", zap.Error(err))
		return err
	}
	user, err := mysql.GetUserByUsername(username)
	if err != nil {
		zap.L().Error("DeleteCommentService() dao.mysql.nzx_sql.GetUserByUsername err=", zap.Error(err))
		return err
	}
	// 只有自己和文章作者和管理员能删除评论
	if username != comment.User.Username && username != comment.Article.User.Username && !user.IsManager {
		fmt.Println(user.IsManager)
		return myErr.OverstepCompetence
	}

	err = mysql.DB.Transaction(func(db *gorm.DB) error {
		if comment.Pid == 0 {
			// 若为一级评论
			// 删除子评论
			result := db.Where("pid = ?", comment.ID).Delete(&gorm_model.Comment{})
			if result.Error != nil {
				zap.L().Error("DeleteComment() dao.mysql.nzx_sql.Delete err=", zap.Error(result.Error))
				return result.Error
			}
			// 删除父级评论
			if err = db.Delete(&comment).Error; err != nil {
				zap.L().Error("DeleteComment() dao.mysql.nzx_sql.Delete err=", zap.Error(err))
				return err
			}

			//	减少文章评论数
			num, err := mysql.QueryArticleCommentNum(int(comment.ArticleID))
			if err != nil {
				zap.L().Error("DeleteComment() dao.mysql.nzx_sql.QueryArticleCommentNum err=", zap.Error(err))
				return err
			}
			err = mysql.UpdateArticleCommentNum(int(comment.ArticleID), num-int(result.RowsAffected)-1, db)
			if err != nil {
				zap.L().Error("DeleteComment() dao.mysql.nzx_sql.UpdateArticleCommentNum err=", zap.Error(err))
				return err
			}

		} else {
			// 若为二级评论
			if err = mysql.DB.Where("id = ?", comment.ID).Delete(&gorm_model.Comment{}).Error; err != nil {
				zap.L().Error("DeleteComment() dao.mysql.nzx_sql.Delete err=", zap.Error(err))
				return err
			}

			//	减少文章评论数
			num, err := mysql.QueryArticleCommentNum(int(comment.ArticleID))
			if err != nil {
				zap.L().Error("DeleteComment() dao.mysql.nzx_sql.QueryArticleCommentNum err=", zap.Error(err))
				return err
			}
			err = mysql.UpdateArticleCommentNum(int(comment.ArticleID), num-1, db)
			if err != nil {
				zap.L().Error("DeleteComment() dao.mysql.nzx_sql.UpdateArticleCommentNum err=", zap.Error(err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		zap.L().Error("DeleteComment() dao.mysql.nzx_sql.Transaction err=", zap.Error(err))
		return err
	}

	// 删除redis
	redis.RDB.HDel("comment", strconv.Itoa(cid))
	return nil
}
