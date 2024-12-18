package mysql

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"studentGrow/models/gorm_model"
)

// UpdateCollectNum 修改收藏数量
func UpdateCollectNum(aid, collectNum int, db *gorm.DB) error {
	if err := db.Model(&gorm_model.Article{}).Where("id = ?", aid).Update("collect_amount", collectNum).Error; err != nil {
		zap.L().Error("UpdateCollectNum() dao.mysql.mysql_collect.Update err=", zap.Error(err))
		return err
	}

	return nil
}

// QueryCollectNum 查询收藏数量
func QueryCollectNum(aid int) (int, error) {
	article := gorm_model.Article{}
	if err := DB.Model(&gorm_model.Article{}).Where("id = ?", aid).First(&article).Error; err != nil {
		zap.L().Error("QueryCollectNum() dao.mysql.mysql_collect.First err=", zap.Error(err))
		return -1, err
	}

	return article.CollectAmount, nil
}

// InsertCollectRecord 插入收藏记录
func InsertCollectRecord(aid, uid int, db *gorm.DB) error {
	if err := db.Model(&gorm_model.UserCollectRecord{}).Create(&gorm_model.UserCollectRecord{UserID: uint(uid), ArticleID: uint(aid)}).Error; err != nil {
		zap.L().Error("InsertCollectRecord() dao.mysql.mysql_collect.Create err=", zap.Error(err))
		return err
	}
	return nil
}

// DeleteCollectRecord 删除收藏记录
func DeleteCollectRecord(aid, uid int, db *gorm.DB) error {
	if err := db.Where("article_id = ? and user_id = ?", aid, uid).Delete(&gorm_model.UserCollectRecord{}).Error; err != nil {
		zap.L().Error("DeleteCollectRecord() dao.mysql.mysql_collect.Delete err=", zap.Error(err))
		return err
	}
	return nil
}
