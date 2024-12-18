package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"studentGrow/models/gorm_model"
	myErr "studentGrow/pkg/error"
	time "studentGrow/utils/timeConverter"
)

// GetUserByUsername 通过username获取user对象
func GetUserByUsername(username string) (*gorm_model.User, error) {
	var user gorm_model.User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		zap.L().Error("GetClassByUsername() dao.mysql.sql_msg.First err=", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

// GetUnreadReportsForClass 获取未读举报信息-班级
func GetUnreadReportsForClass(username string, limit, page int) ([]gorm_model.UserReportArticleRecord, error) {
	// 通过username获取管理员
	user, err := GetUserByUsername(username)
	if err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.GetUserByUsername err=", zap.Error(err))
		return nil, err
	}

	// 查询该班级的所有成员id
	var uids []int64
	if err = DB.Model(&gorm_model.User{}).Where("class = ?", user.Class).Pluck("id", &uids).Error; err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return nil, err
	}

	// 筛选属于该班级的文章
	var aids []int64
	if err = DB.Model(&gorm_model.Article{}).Where("ban = ? AND user_id IN ?", false, uids).Pluck("id", &aids).Error; err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return nil, err
	}

	// 按举报时间逆序查询
	//  通过文章id查找到到对应的用户
	var reports []gorm_model.UserReportArticleRecord
	if err = DB.Model(&gorm_model.UserReportArticleRecord{}).Preload("User").Preload("Article").
		Where("article_id IN ? AND is_read = ?", aids, false).Order("created_at DESC, article_id ASC").
		Limit(limit).Offset((page - 1) * limit).
		Find(&reports).Error; err != nil {
		zap.L().Error("GetClassByUsername() dao.mysql.sql_msg.Find err=", zap.Error(err))
		return nil, err
	}

	return reports, nil
}

// GetUnreadReportNumForClass 获取未读举报信息数目-班级
func GetUnreadReportNumForClass(username string) (int, error) {
	// 通过username获取管理员
	user, err := GetUserByUsername(username)
	if err != nil {
		zap.L().Error("GetUnreadReportNumForClass() dao.mysql.sql_msg.GetUserByUsername err=", zap.Error(err))
		return -1, err
	}

	// 查询该班级的所有成员id
	var uids []int64
	if err = DB.Model(&gorm_model.User{}).Where("class = ?", user.Class).Pluck("id", &uids).Error; err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return -1, err
	}

	// 筛选文章
	var aids []int64
	if err = DB.Model(&gorm_model.Article{}).Where("ban = ? AND user_id IN ?", false, uids).Pluck("id", &aids).Error; err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return -1, err
	}

	// 按举报时间逆序查询
	// 通过文章id查找到到对应的用户
	var count int64
	if err = DB.Model(&gorm_model.UserReportArticleRecord{}).
		Where("is_read = ? AND article_id IN ?", false, aids).
		Count(&count).Error; err != nil {
		zap.L().Error("GetClassByUsername() dao.mysql.sql_msg.Find err=", zap.Error(err))
		return -1, err
	}

	return int(count), nil
}

// GetUnreadReportsForGrade 获取未读举报信息-年级
func GetUnreadReportsForGrade(grade int, limit, page int) ([]gorm_model.UserReportArticleRecord, error) {
	// 通过年级获取入学年份
	year, err := time.GetEnrollmentYear(grade)
	if err != nil {
		zap.L().Error("GetUnreadReportsForGrade() dao.mysql.sql_msg.myErr.GetEnrollmentYear() err=", zap.Error(err))
		return nil, err
	}

	// 查询该年级的所有成员id
	var uids []int64
	if err = DB.Model(&gorm_model.User{}).Where("plus_time between ? and ?", fmt.Sprintf("%d-01-01", year.Year()), fmt.Sprintf("%d-12-31", year.Year())).Pluck("id", &uids).Error; err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return nil, err
	}

	// 筛选该年级文章
	var aids []int64
	if err = DB.Model(&gorm_model.Article{}).Where("ban = ? AND user_id IN ?", false, uids).Pluck("id", &aids).Error; err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return nil, err
	}

	// 按举报时间逆序查询
	// 通过文章id查找到到对应的用户
	var reports []gorm_model.UserReportArticleRecord
	if err = DB.Preload("User").Preload("Article").
		Where("is_read = ? AND article_id IN ?", false, aids).Order("created_at DESC, article_id ASC").
		Limit(limit).Offset((page - 1) * limit).
		Find(&reports).Error; err != nil {
		zap.L().Error("GetUnreadReportsForGrade() dao.mysql.sql_msg.myErr.Find() err=", zap.Error(err))
		return nil, err
	}

	return reports, nil
}

// GetUnreadReportNumForGrade 获取未读举报信息数目-年级
func GetUnreadReportNumForGrade(grade int) (int, error) {
	// 通过年级获取入学年份
	year, err := time.GetEnrollmentYear(grade)
	if err != nil {
		zap.L().Error("GetUnreadReportsForGrade() dao.mysql.sql_msg.myErr.GetEnrollmentYear() err=", zap.Error(err))
		return -1, err
	}

	// 查询该年级的所有成员id
	var uids []int64
	if err = DB.Model(&gorm_model.User{}).Where("plus_time between ? and ?", fmt.Sprintf("%d-01-01", year.Year()), fmt.Sprintf("%d-12-31", year.Year())).Pluck("id", &uids).Error; err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return -1, err
	}

	// 筛选该年级文章
	var aids []int64
	if err = DB.Model(&gorm_model.Article{}).Where("ban = ? AND user_id IN ?", false, uids).Pluck("id", &aids).Error; err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return -1, err
	}

	var count int64
	// 按举报时间逆序查询
	// 通过文章id查找到到对应的用户
	if err = DB.Model(&gorm_model.UserReportArticleRecord{}).Preload("User").Preload("Article").
		Where("is_read = ? AND article_id IN ?", false, aids).
		Count(&count).Error; err != nil {
		zap.L().Error("GetUnreadReportsForGrade() dao.mysql.sql_msg.myErr.Find() err=", zap.Error(err))
		return -1, err
	}

	return int(count), nil
}

// GetUnreadReportsForSuperman 获取未读举报信息 - 超级(院级)
func GetUnreadReportsForSuperman(limit, page int) ([]gorm_model.UserReportArticleRecord, error) {

	var reports []gorm_model.UserReportArticleRecord
	fmt.Println("limit", limit, "page", page)
	if err := DB.Joins("JOIN articles ON user_report_article_records.article_id = articles.id AND articles.ban = ? AND articles.deleted_at IS NULL", false).
		Where("user_report_article_records.is_read = ?", false).
		Order("user_report_article_records.created_at DESC, user_report_article_records.article_id ASC").
		Limit(limit).Offset((page - 1) * limit).
		Preload("User").Preload("Article").
		Find(&reports).Error; err != nil {
		zap.L().Error("GetUnreadReportsForSuperman() dao.mysql.sql_msg.myErr.Find() err=", zap.Error(err))
		return nil, err
	}

	return reports, nil
}

// GetUnreadReportNumForSuperman 获取未读举报信息数目 - 超级(院级)
func GetUnreadReportNumForSuperman() (int, error) {
	var count int64

	if err := DB.Model(&gorm_model.UserReportArticleRecord{}).Joins("JOIN articles ON user_report_article_records.article_id = articles.id AND articles.ban = ?", false).
		Where("user_report_article_records.is_read = ?", false).
		Order("user_report_article_records.created_at DESC, user_report_article_records.article_id ASC").
		Preload("User").Preload("Article").
		Count(&count).Error; err != nil {
		zap.L().Error("GetUnreadReportsForSuperman() dao.mysql.sql_msg.myErr.Find() err=", zap.Error(err))
		return -1, err
	}

	return int(count), nil

}

// AckUnreadReportsForClass 确认未读举报信息 - 班级
func AckUnreadReportsForClass(reportsId int, username string) error {
	// 检查权限(查看所确认文章是否是该班级的文章)
	user, err := GetUserByUsername(username)
	if err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.GetUserByUsername err=", zap.Error(err))
		return err
	}
	// 查询该班级的所有成员id
	var uids []int64
	if err = DB.Model(&gorm_model.User{}).Where("class = ?", user.Class).Pluck("id", &uids).Error; err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return err
	}
	var count int64
	if err = DB.Model(&gorm_model.Article{}).Where("ban = ? AND user_id IN ?", false, uids).Count(&count).Error; err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return err
	}

	if count <= 0 {
		return myErr.OverstepCompetence
	}

	// 修改举报信息读取状态为已读
	result := DB.Where("article_id = ?", reportsId).
		Updates(gorm_model.UserReportArticleRecord{IsRead: true})

	if result.Error != nil {
		zap.L().Error("AckUnreadReportsForClass() dao.mysql.sql_msg.myErr.Find() err=", zap.Error(result.Error))
		return result.Error
	}

	return nil
}

// AckUnreadReportsForGrade 确认未读举报信息 - 年级
func AckUnreadReportsForGrade(reportsId int, grade int) error {
	// 通过年级获取入学年份
	year, err := time.GetEnrollmentYear(grade)
	if err != nil {
		zap.L().Error("AckUnreadReportsForGrade() dao.mysql.sql_msg.myErr.Find() err=", zap.Error(err))
		return err
	}

	// 检查权限(查看所确认文章是否是该年级的文章)
	// 查询该年级的所有成员id
	var uids []int64
	if err = DB.Model(&gorm_model.User{}).Where("plus_time between ? and ?", fmt.Sprintf("%d-01-01", year.Year()), fmt.Sprintf("%d-12-31", year.Year())).Pluck("id", &uids).Error; err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return err
	}
	var count int64
	if err = DB.Model(&gorm_model.Article{}).Where("ban = ? AND user_id IN ?", false, uids).Count(&count).Error; err != nil {
		zap.L().Error("GetUnreadReportsForClass() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return err
	}

	if count <= 0 {
		return myErr.OverstepCompetence
	}

	// 修改举报信息读取状态为已读
	result := DB.Where("article_id = ?", reportsId).
		Updates(gorm_model.UserReportArticleRecord{IsRead: true})

	if result.Error != nil {
		zap.L().Error("AckUnreadReportsForGrade() dao.mysql.sql_msg.myErr.Find() err=", zap.Error(result.Error))
		return result.Error
	}

	return nil
}

// AckUnreadReportsForSuperman 确认未读举报信息 - 超级(院级)
func AckUnreadReportsForSuperman(reportsId int) error {
	// 修改举报信息读取状态为已读
	result := DB.Where("article_id = ?", reportsId).
		Updates(gorm_model.UserReportArticleRecord{IsRead: true})

	if result.Error != nil {
		zap.L().Error("AckUnreadReportsForSuperman() dao.mysql.sql_msg.myErr.Updates() err=", zap.Error(result.Error))
		return result.Error
	}

	return nil
}

// QuerySystemNotification 查询系统消息
func QuerySystemNotification(page, size, uid int) ([]gorm_model.SysNotification, error) {
	var msgs []gorm_model.SysNotification

	if err := DB.Preload("OwnUser").Where("notice_type = ? AND status = ?", 4, false).
		Order("created_at desc").Offset((page - 1) * size).Limit(size).
		Find(&msgs).
		Error; err != nil {
		zap.L().Error("QuerySystemMsg() dao.mysql.sql_msg", zap.Error(err))
		return nil, err
	}
	return msgs, nil
}

// QuerySystemNotificationNum 查询系统消息总数
func QuerySystemNotificationNum() (int, error) {
	var count int64
	if err := DB.Model(&gorm_model.SysNotification{}).Where("notice_type = ? AND status = ?", 4, false).
		Count(&count).Error; err != nil {
		zap.L().Error("QuerySystemNotificationNum() dao.mysql.sql_msg", zap.Error(err))
		return -1, err
	}
	return int(count), nil
}

// QuerySystemNotificationIds 查询所有普通系统消息id
func QuerySystemNotificationIds() ([]int, error) {
	var ids []int
	if err := DB.Model(&gorm_model.SysNotification{}).Where("notice_type = ? AND status = ?", 4, false).
		Pluck("id", &ids).Error; err != nil {
		zap.L().Error("QuerySystemNotificationNum() dao.mysql.sql_msg", zap.Error(err))
		return nil, err
	}
	return ids, nil
}

// QueryArticleBanNotification 查询封禁消息
func QueryArticleBanNotification(uid int) ([]gorm_model.SysNotification, error) {
	var msgs []gorm_model.SysNotification
	if err := DB.Model(&gorm_model.SysNotification{}).
		Where("notice_type = ? AND tar_user_id = ? AND status = ?", 5, uid, false).
		Find(&msgs).Error; err != nil {
		zap.L().Error("QueryArticleBanNotification() dao.mysql.sql_msg", zap.Error(err))
		return nil, err
	}
	return msgs, nil
}

// QueryArticleBanNotificationNum 查询封禁消息总数
func QueryArticleBanNotificationNum(uid int) (int, error) {
	var count int64
	if err := DB.Model(&gorm_model.SysNotification{}).Where("notice_type = ? AND tar_user_id = ? AND status = ?", 5, uid, false).
		Count(&count).Error; err != nil {
		zap.L().Error("QueryArticleBanNotificationNum() dao.mysql.sql_msg", zap.Error(err))
		return -1, err
	}
	return int(count), nil
}

// QueryArticleBanNotificationIds 查询所有封禁消息id
func QueryArticleBanNotificationIds(uid int) ([]int, error) {
	var ids []int
	if err := DB.Model(&gorm_model.SysNotification{}).
		Where("notice_type = ? AND tar_user_id = ? AND status = ?", 5, uid, false).
		Pluck("id", &ids).Error; err != nil {
		zap.L().Error("QueryArticleBanNotificationIds() dao.mysql.sql_msg", zap.Error(err))
		return nil, err
	}
	return ids, nil
}

// QueryManagerNotification 查询管理员消息
func QueryManagerNotification(page, size int) ([]gorm_model.SysNotification, error) {
	var msgs []gorm_model.SysNotification
	if err := DB.Preload("OwnUser").Where("notice_type = ? AND status = ?", 3, false).
		Order("created_at desc").Offset((page - 1) * size).Limit(size).
		Find(&msgs).Error; err != nil {
		zap.L().Error("QueryManagerNotification() dao.mysql.sql_msg", zap.Error(err))
		return nil, err
	}
	return msgs, nil
}

// QueryManagerNotificationNum 查询管理员消息总数
func QueryManagerNotificationNum() (int, error) {
	var count int64
	if err := DB.Model(&gorm_model.SysNotification{}).Where("notice_type = ? AND status = ?", 3, false).Count(&count).Error; err != nil {
		zap.L().Error("QueryManagerNotificationNum() dao.mysql.sql_msg", zap.Error(err))
		return -1, err
	}
	return int(count), nil
}

// QueryManagerNotificationIds 查询所有管理员消息id
func QueryManagerNotificationIds() ([]int, error) {
	var ids []int
	if err := DB.Model(&gorm_model.SysNotification{}).Where("notice_type = ? AND status = ?", 3, false).Pluck("id", &ids).Error; err != nil {
		zap.L().Error("QueryManagerNotificationIds() dao.mysql.sql_msg", zap.Error(err))
		return nil, err
	}
	return ids, nil
}

// QueryLikeRecordByUser 分页查询其文章和评论的点赞记录
func QueryLikeRecordByUser(uid, page, limit int) ([]gorm_model.UserLikeRecord, error) {
	var likes []gorm_model.UserLikeRecord
	if err := DB.Preload("User").Preload("Article").Preload("Comment.Article").
		Where("(article_id IN (SELECT id FROM articles WHERE user_id = ? AND ban = ? AND status = ? AND deleted_at IS NULL) OR comment_id IN (SELECT id FROM comments WHERE user_id = ? AND deleted_at IS NULL))", uid, false, true, uid).
		Limit(limit).
		Offset((page - 1) * limit).Order("created_at desc").
		Find(&likes).Error; err != nil {
		zap.L().Error("QueryLikeRecordByUser() dao.mysql.sql_msg.Find err=", zap.Error(err))
		return nil, err
	}

	return likes, nil
}

// QueryLikeRecordNumByUser 查询未读点赞记录数量
func QueryLikeRecordNumByUser(uid int) (int, error) {
	var count int64

	if err := DB.Model(&gorm_model.UserLikeRecord{}).Where("is_read = ? AND (article_id IN (SELECT id FROM articles WHERE user_id = ? AND ban = ? AND status = ? AND deleted_at IS NULL) OR comment_id IN (SELECT id FROM comments WHERE user_id = ? AND deleted_at IS NULL))", false, uid, false, true, uid).Count(&count).Error; err != nil {
		zap.L().Error("QueryLikeRecordNumByUserArticle() dao.mysql.sql_msg.Count err=", zap.Error(err))
		return -1, err
	}
	return int(count), nil
}

// QueryCollectRecordByUserArticles 通过用户的所有文章查找其收藏记录(该用户的文章被谁收藏了记录)
func QueryCollectRecordByUserArticles(uid, page, limit int) ([]gorm_model.UserCollectRecord, error) {
	// 获取该用户文章列表
	aids, err := QueryArticleIdsByUserId(uid)
	if err != nil {
		zap.L().Error("QueryCollectRecordByUserArticles() dao.mysql.sql_msg.QueryArticleIdsByUserId err=", zap.Error(err))
		return nil, err
	}

	// 通过文章id查询收藏记录
	var articleCollects []gorm_model.UserCollectRecord
	if err := DB.Preload("Article").Preload("User").Where("article_id IN ?", aids).
		Order("created_at desc").
		Limit(limit).Offset((page - 1) * limit).
		Find(&articleCollects).Error; err != nil {
		zap.L().Error("QueryCollectRecordByUserArticles() dao.mysql.sql_msg.Find err=", zap.Error(err))
		return nil, err
	}

	return articleCollects, nil
}

// QueryCollectRecordNumByUserArticle 通过uid查询其文章的未读收藏记录数量
func QueryCollectRecordNumByUserArticle(uid int) (int, error) {
	// 获取该用户文章列表
	aids, err := QueryArticleIdsByUserId(uid)
	if err != nil {
		zap.L().Error("QueryCollectRecordByUserArticles() dao.mysql.sql_msg.QueryArticleIdsByUserId err=", zap.Error(err))
		return -1, err
	}

	var count int64

	if err := DB.Model(&gorm_model.UserCollectRecord{}).Where("article_id IN ? and is_read = ?", aids, false).Count(&count).Error; err != nil {
		zap.L().Error("QueryCollectRecordNumByUserArticle() dao.mysql.sql_msg.Count err=", zap.Error(err))
		return -1, err
	}

	return int(count), nil
}

// QueryCommentRecordByUserArticles 通过用户的所有文章和评论查找其评论记录(该用户的文章或评论被谁评论了记录)
func QueryCommentRecordByUserArticles(uid, page, limit int) (gorm_model.Comments, error) {
	var comments gorm_model.Comments
	var commentIDs []int

	if err := DB.Model(&gorm_model.Comment{}).Where("user_id = ? AND pid = ?", uid, 0).Pluck("id", &commentIDs).Error; err != nil {
		zap.L().Error("QueryCommentRecordByUserArticles() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return nil, err
	}

	var articleIDs []int
	if err := DB.Model(&gorm_model.Article{}).Where("user_id = ? AND ban = ? AND status = ?", uid, false, true).Pluck("id", &articleIDs).Error; err != nil {
		zap.L().Error("QueryCommentRecordByUserArticles() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return nil, err
	}

	// 查找到的是回复评论的评论内容以及评论文章的评论内容
	if err := DB.Model(&gorm_model.Comment{}).Preload("User").Preload("Article").Limit(limit).Offset((page-1)*limit).
		Order("created_at desc").
		Where("pid IN ?", commentIDs).Or("article_id IN ? AND pid = ?", articleIDs, 0).Find(&comments).Error; err != nil {
		zap.L().Error("QueryCommentRecordByUserArticles() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return nil, err
	}

	return comments, nil
}

// QueryCommentRecordNumByUserId 通过uid获取未读评论记录数量
func QueryCommentRecordNumByUserId(uid int) (int, error) {
	var count int64

	var commentIDs []int

	if err := DB.Model(&gorm_model.Comment{}).Where("user_id = ? AND pid = ?", uid, 0).Pluck("id", &commentIDs).Error; err != nil {
		zap.L().Error("QueryCommentRecordByUserArticles() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return -1, err
	}

	var articleIDs []int
	if err := DB.Model(&gorm_model.Article{}).Where("user_id = ? AND ban = ? AND status = ?", uid, false, true).Pluck("id", &articleIDs).Error; err != nil {
		zap.L().Error("QueryCommentRecordByUserArticles() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return -1, err
	}

	// 查找到的是回复评论的评论内容以及评论文章的评论内容
	if err := DB.Model(&gorm_model.Comment{}).Preload("User").Preload("Article").
		Where("pid IN ? AND is_read = ?", commentIDs, false).Or("article_id IN ? AND pid = ? AND is_read = ?", articleIDs, 0, false).Count(&count).Error; err != nil {
		zap.L().Error("QueryCommentRecordByUserArticles() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return -1, err
	}
	return int(count), nil

}

// QueryCommentRecordNumByUserArticle 通过uid查找文章评论未读记录
func QueryCommentRecordNumByUserArticle(uid int) (int, error) {
	var count int64
	if err := DB.Model(gorm_model.Comment{}).Preload("Article", "user_id = ? and ban = ?", uid, false).Where("is_read = ?", false).Count(&count).Error; err != nil {
		zap.L().Error("QueryCommentRecordNumByUserArticle() dao.mysql.sql_msg.Count err=", zap.Error(err))
		return -1, err
	}
	return int(count), nil
}

// QueryCommentRecordByUserComments 通过用户的评论查找其被评论记录
func QueryCommentRecordByUserComments(cid int) (gorm_model.Comments, error) {
	comments := gorm_model.Comments{}
	if err := DB.Where("pid = ?", cid).Order("created_at desc").Find(&comments).Error; err != nil {
		zap.L().Error("QueryCommentRecordByUserComments() dao.mysql.sql_msg.Find err=", zap.Error(err))
		return nil, err
	}

	if len(comments) == 0 {
		zap.L().Error("QueryCommentRecordByUserComments() dao.mysql.sql_msg err=", zap.Error(myErr.ErrNotFoundError))
		return nil, myErr.ErrNotFoundError
	}
	return comments, nil
}

// UpdateLikeRecordRead 确认点赞
func UpdateLikeRecordRead(msgId int) error {
	if err := DB.Model(&gorm_model.UserLikeRecord{}).Where("id = ?", msgId).Update("is_read", true).Error; err != nil {
		zap.L().Error("UpdateArticleLikeRecordRead() dao.mysql.sql_msg err=", zap.Error(err))
		return err
	}
	return nil
}

// UpdateCollectRecordRead 确认收藏
func UpdateCollectRecordRead(msgId int) error {
	if err := DB.Model(&gorm_model.UserCollectRecord{}).Where("id = ?", msgId).Update("is_read", true).Error; err != nil {
		zap.L().Error("UpdateArticleLikeRecordRead() dao.mysql.sql_msg err=", zap.Error(err))
		return err
	}
	return nil
}

// UpdateCommentRecordRead 确认评论
func UpdateCommentRecordRead(cid int) error {
	if err := DB.Model(&gorm_model.Comment{}).Where("id = ?", cid).Update("is_read", true).Error; err != nil {
		zap.L().Error("UpdateArticleLikeRecordRead() dao.mysql.sql_msg err=", zap.Error(err))
		return err
	}
	return nil
}

// AddManagerNotification 添加管理员通知
func AddManagerNotification(ownId int, content string) error {
	//managerMsg := gorm_model.MsgRecord{
	//	Content:  content,
	//	Username: username,
	//	UserID:   uint(uid),
	//	Type:     constant.ManagerMsgConstant,
	//}
	//
	//if err := DB.Create(&managerMsg).Error; err != nil {
	//	zap.L().Error("AddManagerMsg() dao.mysql.sql_msg err=", zap.Error(err))
	//	return err
	//}

	sysNotification := gorm_model.SysNotification{
		OwnUserId:  uint(ownId),
		NoticeType: 3,
		Content:    content,
	}

	if err := DB.Create(&sysNotification).Error; err != nil {
		zap.L().Error("AddManagerNotification() dao.mysql.sql_msg err=", zap.Error(err))
		return err
	}

	return nil
}

// AddSystemNotification 添加普通系统通知
func AddSystemNotification(content string, ownId int) error {
	//systemMsg := gorm_model.MsgRecord{
	//	Username: username,
	//	Content:  content,
	//	UserID:   uint(uid),
	//	Type:     constant.SystemMsgConstant,
	//}
	//
	//if err := db.Create(&systemMsg).Error; err != nil {
	//	zap.L().Error("AddSystemMsg() dao.mysql.sql_msg err=", zap.Error(err))
	//	return err
	//}

	sysNotification := gorm_model.SysNotification{
		OwnUserId:  uint(ownId),
		NoticeType: 4,
		Content:    content,
	}

	if err := DB.Create(&sysNotification).Error; err != nil {
		zap.L().Error("AddSystemNotification() dao.mysql.sql_msg err=", zap.Error(err))
		return err
	}

	return nil
}

// AddBanSystemNotification 添加文章封禁通知
func AddBanSystemNotification(content string, ownId int, db *gorm.DB, tarId int) error {
	sysNotification := gorm_model.SysNotification{
		OwnUserId:  uint(ownId),
		TarUserId:  uint(tarId),
		NoticeType: 5,
		Content:    content,
	}

	if err := db.Create(&sysNotification).Error; err != nil {
		zap.L().Error("AddBanSystemNotification() dao.mysql.sql_msg err=", zap.Error(err))
		return err
	}

	return nil
}

// DeleteSystemNotification 删除系统消息
func DeleteSystemNotification(MsgId int, db *gorm.DB) error {
	if err := db.Model(&gorm_model.SysNotification{}).Delete("notice_type = ? AND id = ?", 4, MsgId).Update("status", true).Error; err != nil {
		zap.L().Error("DeleteSystemNotification() dao.mysql.sql_msg err=", zap.Error(err))
		return err
	}
	return nil
}

// DeleteManagerNotification 删除管理员消息
func DeleteManagerNotification(MsgId int, db *gorm.DB) error {
	if err := db.Model(&gorm_model.SysNotification{}).Where("notice_type = ? AND id = ?", 3, MsgId).Update("status", true).Error; err != nil {
		zap.L().Error("DeleteSystemNotification() dao.mysql.sql_msg err=", zap.Error(err))
		return err
	}
	return nil
}

// AckUserAllLikeId 确认用户所有未读的被点赞消息id
func AckUserAllLikeId(uid int) error {
	if err := DB.Model(&gorm_model.UserLikeRecord{}).Where("is_read = ? AND (article_id IN (SELECT id FROM articles WHERE user_id = ? AND ban = ? AND status = ? AND deleted_at IS NULL) OR comment_id IN (SELECT id FROM comments WHERE user_id = ? AND deleted_at IS NULL))", false, uid, false, true, uid).Update("is_read", true).Error; err != nil {
		zap.L().Error("QueryLikeRecordNumByUserArticle() dao.mysql.sql_msg.Count err=", zap.Error(err))
		return err
	}
	return nil
}

// AckUserAllCommentId 确认用户所有未读的被评论消息id
func AckUserAllCommentId(uid int) error {
	var commentIDs []int

	if err := DB.Model(&gorm_model.Comment{}).Where("user_id = ? AND pid = ?", uid, 0).Pluck("id", &commentIDs).Error; err != nil {
		zap.L().Error("QueryUserAllCommentId() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return err
	}

	var articleIDs []int
	if err := DB.Model(&gorm_model.Article{}).Where("user_id = ? AND ban = ? AND status = ?", uid, false, true).Pluck("id", &articleIDs).Error; err != nil {
		zap.L().Error("QueryCommentRecordByUserArticles() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return err
	}

	if err := DB.Model(&gorm_model.Comments{}).
		Where("pid IN ? AND is_read = ?", commentIDs, false).Or("article_id IN ? AND pid = ? AND is_read = ?", articleIDs, 0, false).
		Update("is_read", true).Error; err != nil {
		zap.L().Error("QueryUserAllCommentId() dao.mysql.sql_msg.Pluck err=", zap.Error(err))
		return err
	}
	return nil
}

// AckUserAllCollectId 确认用户所有未读的收藏消息id
func AckUserAllCollectId(uid int) error {
	// 获取该用户文章列表
	aids, err := QueryArticleIdsByUserId(uid)
	if err != nil {
		zap.L().Error("QueryUserAllCollectId() dao.mysql.sql_msg.QueryArticleIdsByUserId err=", zap.Error(err))
		return err
	}

	if err := DB.Model(&gorm_model.UserCollectRecord{}).Where("article_id IN ? and is_read = ?", aids, false).Update("is_read", true).Error; err != nil {
		zap.L().Error("QueryUserAllCollectId() dao.mysql.sql_msg.Count err=", zap.Error(err))
		return err
	}

	return nil
}
