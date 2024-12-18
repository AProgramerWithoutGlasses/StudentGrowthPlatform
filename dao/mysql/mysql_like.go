package mysql

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"studentGrow/models/constant"
	"studentGrow/models/gorm_model"
	myErr "studentGrow/pkg/error"
)

// UpdateLikeNum 修改点赞数量
func UpdateLikeNum(objId, likeType, likeNum int, db *gorm.DB) error {
	switch likeType {
	// 修改文章点赞
	case constant.ArticleInteractionConstant:
		if err := db.Model(gorm_model.Article{}).Where("id = ?", objId).Update("like_amount", likeNum).Error; err != nil {
			zap.L().Error("UpdateLikeNum() dao.mysql.mysql_like.Update err=", zap.Error(err))
			return err
		}
		// 修改评论点赞
	case constant.CommentInteractionConstant:
		if err := db.Model(gorm_model.Comment{}).Where("id = ?", objId).Update("like_amount", likeNum).Error; err != nil {
			zap.L().Error("UpdateLikeNum() dao.mysql.mysql_like.Update err=", zap.Error(err))
			return err
		}
	}
	return nil
}

// QueryLikeNum 查询点赞数量
func QueryLikeNum(objId, likeType int) (int, error) {
	article := gorm_model.Article{}
	comment := gorm_model.Comment{}
	switch likeType {
	case constant.ArticleInteractionConstant:
		if err := DB.Model(gorm_model.Article{}).Where("id = ?", objId).First(&article).Error; err != nil {
			zap.L().Error("QueryLikeNum() dao.mysql.mysql_like.First err=", zap.Error(err))
			return -1, err
		}
		return article.LikeAmount, nil
	case constant.CommentInteractionConstant:
		if err := DB.Model(gorm_model.Comment{}).Where("id = ?", objId).First(&comment).Error; err != nil {
			zap.L().Error("QueryLikeNum() dao.mysql.mysql_like.First err=", zap.Error(err))
			return -1, err
		}
		return comment.LikeAmount, nil
	default:
		return -1, myErr.DataFormatError()
	}
}

// InsertLikeRecord 插入点赞记录
func InsertLikeRecord(objId, likeType int, uid int, db *gorm.DB) error {

	switch likeType {
	case constant.ArticleInteractionConstant:
		articleLike := gorm_model.UserLikeRecord{ArticleID: uint(objId), UserID: uint(uid), Type: likeType}
		if err := db.Model(gorm_model.UserLikeRecord{}).Create(&articleLike).Error; err != nil {
			zap.L().Error("InsertLikeRecord() dao.mysql.mysql_like.Create err=", zap.Error(err))
			return err
		}
	case constant.CommentInteractionConstant:
		commentLike := gorm_model.UserLikeRecord{CommentID: uint(objId), UserID: uint(uid), Type: likeType}
		if err := db.Model(gorm_model.UserLikeRecord{}).Create(&commentLike).Error; err != nil {
			zap.L().Error("InsertLikeRecord() dao.mysql.mysql_like.Create err=", zap.Error(err))
			return err
		}
	default:
		return myErr.DataFormatError()
	}
	return nil
}

// DeleteLikeRecord 删除点赞记录
func DeleteLikeRecord(objId, likeType, uid int) error {
	switch likeType {
	case constant.ArticleInteractionConstant:
		if err := DB.Where("article_id = ? and user_id = ?", objId, uid).Delete(&gorm_model.UserLikeRecord{}).Error; err != nil {
			zap.L().Error("DeleteLikeRecord() dao.mysql.mysql_like.Delete err=", zap.Error(err))
			return err
		}
	case constant.CommentInteractionConstant:
		if err := DB.Where("comment_id = ? and user_id = ?", objId, uid).Delete(&gorm_model.UserLikeRecord{}).Error; err != nil {
			zap.L().Error("DeleteLikeRecord() dao.mysql.mysql_like.Delete err=", zap.Error(err))
			return err
		}
	default:
		return myErr.DataFormatError()
	}
	return nil
}
