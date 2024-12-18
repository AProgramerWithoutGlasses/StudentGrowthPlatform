package message

import (
	"fmt"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	"studentGrow/dao/redis"
	"studentGrow/models/constant"
	"studentGrow/models/gorm_model"
	"studentGrow/models/nzx_model"
	myErr "studentGrow/pkg/error"
	"studentGrow/pkg/sse"
	NotificationPush "studentGrow/service/notificationPush"
	"studentGrow/utils/timeConverter"
)

// GetSystemMsgService 获取系统消息通知
func GetSystemMsgService(limit, page int, username string) ([]gorm_model.SysNotification, int, error) {

	uid, err := mysql.QueryUserIdByUsername(username)
	if err != nil {
		zap.L().Error("GetSystemMsgService() service.message.QueryUserIdByUsername", zap.Error(err))
		return nil, 0, err
	}
	// 查询普通系统消息列表
	msgs, err := mysql.QuerySystemNotification(page, limit, uid)
	if err != nil {
		zap.L().Error("GetSystemMsgService() service.message.QuerySystemNotification", zap.Error(err))
		return nil, 0, err
	}
	// 查询封禁消息列表
	msgs2, err := mysql.QueryArticleBanNotification(uid)
	if err != nil {
		zap.L().Error("GetSystemMsgService() service.message.QueryArticleBanNotification", zap.Error(err))
		return nil, 0, err
	}
	msgs = append(msgs, msgs2...)
	for i := 0; i < len(msgs); i++ {
		msgs[i].Time = timeConverter.IntervalConversion(msgs[i].CreatedAt)

		isAck, err := redis.IsUserAckedSystemNotification(uid, int(msgs[i].ID))
		if err != nil {
			zap.L().Error("GetSystemMsgService() service.message.IsUserAckedSystemNotification", zap.Error(err))
			return nil, 0, err
		}
		if isAck {
			msgs[i].IsRead = true
		} else {
			msgs[i].IsRead = false
		}
	}

	// 查询系统消息总数
	total, err := mysql.QuerySystemNotificationNum()
	if err != nil {
		zap.L().Error("GetSystemMsgService() service.message.QuerySystemNotificationNum", zap.Error(err))
		return nil, 0, err
	}

	// 查询封禁消息总数
	total2, err := mysql.QueryArticleBanNotificationNum(uid)
	if err != nil {
		zap.L().Error("GetSystemMsgService() service.message.QueryArticleBanNotificationNum", zap.Error(err))
		return nil, 0, err
	}
	total += total2
	// 查询用户确认信息数
	ackNum, err := redis.GetUserAckedSystemNum(uid)
	if err != nil {
		zap.L().Error("GetSystemMsgService() service.message.GetUserAckedNum", zap.Error(err))
		return nil, 0, err
	}

	count := total - ackNum
	fmt.Println("total", total, "ackNum", ackNum)
	return msgs, count, nil
}

// GetManagerMsgService 获取管理员消息
func GetManagerMsgService(limit, page int, username string) ([]gorm_model.SysNotification, int, error) {
	uid, err := mysql.QueryUserIdByUsername(username)
	if err != nil {
		zap.L().Error("GetManagerMsgService() service.message.QueryUserIdByUsername", zap.Error(err))
		return nil, 0, err
	}

	// 查询管理员消息列表
	msgs, err := mysql.QueryManagerNotification(page, limit)
	if err != nil {
		zap.L().Error("GetManagerMsgService() service.message.QueryManagerNotification", zap.Error(err))
		return nil, 0, err
	}

	for i := 0; i < len(msgs); i++ {
		msgs[i].Time = timeConverter.IntervalConversion(msgs[i].CreatedAt)

		isAck, err := redis.IsUserAckedManagerNotification(uid, int(msgs[i].ID))
		if err != nil {
			zap.L().Error("GetSystemMsgService() service.message.IsUserAckedManagerNotification", zap.Error(err))
			return nil, 0, err
		}
		if isAck {
			msgs[i].IsRead = true
		} else {
			msgs[i].IsRead = false
		}
	}

	total, err := mysql.QueryManagerNotificationNum()
	if err != nil {
		zap.L().Error("GetManagerMsgService() service.message.QueryManagerNotificationNum", zap.Error(err))
		return nil, 0, err
	}

	ackNum, err := redis.GetUserAckedManagerNum(uid)
	if err != nil {
		zap.L().Error("GetManagerMsgService() service.message.GetUserAckedManagerNum", zap.Error(err))
		return nil, 0, err
	}

	count := total - ackNum

	return msgs, count, nil
}

// GetArticleAndCommentLikedMsgService  获取点赞消息
func GetArticleAndCommentLikedMsgService(username string, page, limit int) ([]nzx_model.Out, int, error) {
	// 获取uid
	uid, err := mysql.GetIdByUsername(username)
	if err != nil {
		zap.L().Error("GetArticleAndCommentLikedMsgService() service.article.likeService.GetIdByUsername err=", zap.Error(err))
		return nil, -1, err
	}

	// 获取点赞列表
	likes, err := mysql.QueryLikeRecordByUser(uid, page, limit)
	if err != nil {
		zap.L().Error("GetArticleAndCommentLikedMsgService() service.article.likeService.QueryLikeRecordByUserArticle err=", zap.Error(err))
		return nil, -1, err
	}

	// 获取文章点赞未读消息总数
	sum, err := mysql.QueryLikeRecordNumByUser(uid)
	if err != nil {
		zap.L().Error("GetArticleAndCommentLikedMsgService() service.article.likeService.QueryLikeRecordNumByUser err=", zap.Error(err))
		return nil, -1, err
	}

	list := make([]nzx_model.Out, 0)

	for _, like := range likes {
		// 判断文章点赞还是评论点赞
		usernameL := like.User.Username
		name := like.User.Name
		content := like.Article.Content
		userHeadshot := like.User.HeadShot
		likeType := 0
		articleId := like.ArticleID
		if like.Type == 1 {
			content = like.Comment.Content
			likeType = 1
			articleId = like.Comment.ArticleID
		}

		list = append(list, nzx_model.Out{
			Username:     usernameL,
			Name:         name,
			Content:      content,
			UserHeadshot: userHeadshot,
			PostTime:     timeConverter.IntervalConversion(like.CreatedAt),
			IsRead:       like.IsRead,
			Type:         likeType,
			ArticleId:    articleId,
			MsgId:        like.ID,
		})
	}

	return list, sum, nil
}

// GetCollectMsgService 获取收藏消息
func GetCollectMsgService(username string, page, limit int) ([]map[string]any, int, error) {
	// 获取uid
	uid, err := mysql.GetIdByUsername(username)
	if err != nil {
		zap.L().Error("GetCollectMsgService() service.article.likeService.GetIdByUsername err=", zap.Error(err))
		return nil, -1, err
	}

	// 获取收藏消息列表
	articleCollects, err := mysql.QueryCollectRecordByUserArticles(uid, page, limit)
	if err != nil {
		zap.L().Error("GetCollectMsgService() service.article.likeService.QueryCollectRecordByUserArticles err=", zap.Error(err))
		return nil, -1, err
	}

	// 获取未读收藏消息数量
	collectNum, err := mysql.QueryCollectRecordNumByUserArticle(uid)
	if err != nil {
		zap.L().Error("GetCollectMsgService() service.article.likeService.QueryCollectRecordNumByUserArticle err=", zap.Error(err))
		return nil, -1, err
	}

	list := make([]map[string]any, 0)

	for _, collect := range articleCollects {
		list = append(list, map[string]any{
			"username":        collect.User.Username,
			"name":            collect.User.Name,
			"article_content": collect.Article.Content,
			"user_headshot":   collect.User.HeadShot,
			"post_time":       timeConverter.IntervalConversion(collect.CreatedAt),
			"is_read":         collect.IsRead,
			"article_id":      collect.ArticleID,
			"msg_id":          collect.ID,
		})
	}

	return list, collectNum, nil

}

// GetCommentMsgService 获取评论消息
func GetCommentMsgService(username string, page, limit int) (nzx_model.CommentMsgs, int, error) {
	// 获取uid
	uid, err := mysql.GetIdByUsername(username)
	if err != nil {
		zap.L().Error("GetCommentMsgService() service.article.likeService.GetIdByUsername err=", zap.Error(err))
		return nil, -1, err
	}

	// 获取所有评论及回复
	comments, err := mysql.QueryCommentRecordByUserArticles(uid, page, limit)
	if err != nil {
		zap.L().Error("GetCommentMsgService() service.article.likeService.QueryCommentRecordByUserArticles err=", zap.Error(err))
		return nil, -1, err
	}

	commentMsgs := make(nzx_model.CommentMsgs, 0)

	for _, comment := range comments {
		// 默认文章的评论
		commentType := 0
		if comment.Pid != 0 {
			// 如果是评论的回复
			commentType = 1
		}

		commentMsgs = append(commentMsgs, nzx_model.CommentMsg{
			Username:     comment.User.Username,
			Name:         comment.User.Name,
			Content:      comment.Content,
			UserHeadshot: comment.User.HeadShot,
			PostTime:     timeConverter.IntervalConversion(comment.CreatedAt),
			IsRead:       comment.IsRead,
			Type:         commentType,
			ArticleId:    comment.ArticleID,
			MsgId:        comment.ID,
		})
	}

	// 获取未读评论数
	num, err := mysql.QueryCommentRecordNumByUserId(uid)
	if err != nil {
		zap.L().Error("GetCommentMsgService() service.article.likeService.QueryCommentRecordNumByUserId err=", zap.Error(err))
		return nil, -1, err
	}

	return commentMsgs, num, nil
}

// AckInterMsgService 确认互动消息通知
func AckInterMsgService(msgId, msgType int) error {
	switch msgType {
	case constant.LikeMsgConstant:
		err := mysql.UpdateLikeRecordRead(msgId)
		if err != nil {
			zap.L().Error("AckInterMsgService() service.article.likeService.UpdateLikeRecordRead err=", zap.Error(err))
			return err
		}
	case constant.CommentMsgConstant:
		err := mysql.UpdateCommentRecordRead(msgId)
		if err != nil {
			zap.L().Error("AckInterMsgService() service.article.likeService.UpdateCommentRecordRead err=", zap.Error(err))
			return err
		}
	case constant.CollectMsgConstant:
		err := mysql.UpdateCollectRecordRead(msgId)
		if err != nil {
			zap.L().Error("AckInterMsgService() service.article.likeService.UpdateCollectRecordRead err=", zap.Error(err))
			return err
		}
	default:
		return myErr.DataFormatError()
	}
	return nil
}

// AckAllInterMsgService 一键已读互动消息
func AckAllInterMsgService(uid, msgType int) (err error) {
	switch msgType {
	case constant.LikeMsgConstant:
		err = mysql.AckUserAllLikeId(uid)
	case constant.CommentMsgConstant:
		err = mysql.AckUserAllCommentId(uid)
	case constant.CollectMsgConstant:
		err = mysql.AckUserAllCollectId(uid)
	default:
		return myErr.DataFormatError()
	}
	if err != nil {
		zap.L().Error("AckAllInterMsgService() service.article.likeService.AckAllInterMsgService err=", zap.Error(err))
		return err
	}
	return nil
}

// AckManagerMsgService 确认管理员消息
func AckManagerMsgService(username string) error {
	// 获取uid
	uid, err := mysql.GetIdByUsername(username)
	if err != nil {
		zap.L().Error("AckManagerMsgService() service.article.likeService.GetIdByUsername err=", zap.Error(err))
		return err
	}

	// 查询当前所有管理员通知id
	ids, err := mysql.QueryManagerNotificationIds()
	if err != nil {
		zap.L().Error("AckSystemMsgService() service.article.likeService.QueryManagerNotificationIds err=", zap.Error(err))
		return err
	}
	// 将消息加入用户已读消息set集合
	for _, msgId := range ids {
		err = redis.AckManagerNotification(uid, msgId)
		if err != nil {
			zap.L().Error("AckSystemMsgService() service.article.likeService.AckManagerNotification err=", zap.Error(err))
			return err
		}
	}

	// 将用户加入消息set集合
	for _, msgId := range ids {
		err = redis.AddUserToNotificationSet(uid, msgId)
		if err != nil {
			zap.L().Error("AckSystemMsgService() service.article.likeService.AddUserToNotificationHash err=", zap.Error(err))
			return err
		}
	}
	return nil
}

// AckSystemMsgService 确认系统消息
func AckSystemMsgService(username string) error {
	// 获取uid
	uid, err := mysql.GetIdByUsername(username)
	if err != nil {
		zap.L().Error("AckSystemMsgService() service.article.likeService.GetIdByUsername err=", zap.Error(err))
		return err
	}

	// 查询当前所有普通系统通知id
	ids, err := mysql.QuerySystemNotificationIds()
	if err != nil {
		zap.L().Error("AckSystemMsgService() service.article.likeService.QuerySystemNotificationIds err=", zap.Error(err))
		return err
	}
	// 查询当前所有封禁消息id
	ids2, err := mysql.QueryArticleBanNotificationIds(uid)
	if err != nil {
		zap.L().Error("AckSystemMsgService() service.article.likeService.QueryArticleBanNotificationIds err=", zap.Error(err))
		return err
	}
	ids = append(ids, ids2...)
	// 将消息加入用户已读消息set集合
	for _, msgId := range ids {
		err = redis.AckSystemNotification(uid, msgId)
		if err != nil {
			zap.L().Error("AckSystemMsgService() service.article.likeService.GetIdByUsername err=", zap.Error(err))
			return err
		}
	}
	return nil
}

// PublishManagerMsgService 发布管理员通知
func PublishManagerMsgService(username, content, role string) error {
	// 权限验证
	if role != "college" {
		zap.L().Error("PublishManagerMsgService() service.article.likeService.role err=", zap.Error(myErr.OverstepCompetence))
		return myErr.OverstepCompetence
	}

	// 添加管理员通知
	userId, err := mysql.QueryUserIdByUsername(username)
	if err != nil {
		zap.L().Error("PublishManagerMsgService() service.article.QueryUserIdByUsername err=", zap.Error(err))
		return err
	}

	err = mysql.AddManagerNotification(userId, content)
	if err != nil {
		zap.L().Error("PublishManagerMsgService() service.article.AddManagerNotification err=", zap.Error(err))
		return err
	}

	// 消息推送
	notification, err := NotificationPush.BuildManagerNotification(username, content)
	if err != nil {
		zap.L().Error("PublishManagerMsgService() service.article.BuildManagerNotification.Transaction err=", zap.Error(err))
		return err
	}

	sse.SendSysNotification(*notification)

	return nil
}

// PublishSystemMsgService 发布系统通知
func PublishSystemMsgService(content, role, username string) error {
	// 权限验证
	if role != "superman" {
		return myErr.OverstepCompetence
	}

	// 添加通知
	userId, err := mysql.QueryUserIdByUsername(username)
	if err != nil {
		zap.L().Error("PublishSystemMsgService() service.article.QueryUserIdByUsername err=", zap.Error(err))
		return err
	}

	err = mysql.AddSystemNotification(content, userId)
	if err != nil {
		zap.L().Error("PublishSystemMsgService() service.article.AddSystemNotification err=", zap.Error(err))
		return err
	}

	// 消息推送
	notification, err := NotificationPush.BuildSystemNotification(username, content)
	if err != nil {
		zap.L().Error("PublishSystemMsgService() service.article.BuildSystemNotification err=", zap.Error(err))
		return err
	}

	sse.SendSysNotification(*notification)
	return nil
}

// DeleteSystemMsgService 撤销系统消息
func DeleteSystemMsgService(MsgId int) error {

	err := mysql.DeleteSystemNotification(MsgId, mysql.DB)
	if err != nil {
		zap.L().Error("DeleteSystemMsgService() service.article.DeleteSystemNotification err=", zap.Error(err))
		return err
	}
	return nil
}

// DeleteManagerMsgService 撤销管理员消息
func DeleteManagerMsgService(MsgId int, role string) error {
	if role != "college" {
		return myErr.OverstepCompetence
	}
	// 查询已读该消息的用户id
	ids, err := redis.GetUserIdsByNotificationSet(MsgId)
	if err != nil {
		zap.L().Error("DeleteManagerMsgService() service.article.GetUserIdsByNotificationSet err=", zap.Error(err))
		return err
	}

	// 在用户已读set中移除该消息
	for _, id := range ids {
		err = redis.RemoveManagerNotificationInUserAcked(MsgId, id)
		if err != nil {
			zap.L().Error("DeleteManagerMsgService() service.article.RemoveManagerNotificationInUserAcked err=", zap.Error(err))
			return err
		}
	}

	// 移除该消息set
	err = redis.RemoveManagerNotification(MsgId)
	if err != nil {
		zap.L().Error("DeleteManagerMsgService() service.article.RemoveManagerNotification err=", zap.Error(err))
		return err
	}

	// 撤销数据库中的管理员消息
	err = mysql.DeleteManagerNotification(MsgId, mysql.DB)
	if err != nil {
		zap.L().Error("DeleteManagerMsgService() service.article.DeleteManagerNotification err=", zap.Error(err))
		return err
	}
	return nil
}
