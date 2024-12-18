package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"studentGrow/models/gorm_model"
	myErr "studentGrow/pkg/error"
)

// InsertReadRecord 插入浏览记录
func InsertReadRecord(uid, aid int, db *gorm.DB) error {
	readRecord := gorm_model.UserReadRecord{
		UserID:    uint(uid),
		ArticleID: uint(aid),
	}

	if err := db.Create(&readRecord).Error; err != nil {
		zap.L().Error("InsertReadRecord() dao.mysql.mysql_read.Create err=", zap.Error(err))
		return err
	}
	return nil
}

// QueryReadListByUserId 查询浏览记录列表
func QueryReadListByUserId(uid int) (readRecords []gorm_model.UserReadRecord, err error) {
	if err = DB.Where("user_id = ?", uid).Find(&readRecords).Error; err != nil {
		zap.L().Error("QueryReadListByUserId() dao.mysql.sql_article", zap.Error(err))
		return nil, err
	}

	if len(readRecords) == 0 {
		return nil, myErr.ErrNotFoundError
	}
	return readRecords, nil
}

// QueryArticleReadNumById 通过文章id查询文章的浏览量
func QueryArticleReadNumById(aid int) (int, error) {
	var article gorm_model.Article
	if err := DB.Where("id = ?", aid).First(&article).Error; err != nil {
		zap.L().Error("QueryArticleReadNum() dao.mysql.sql_article", zap.Error(err))
		return -1, err
	}
	fmt.Println("readAmount:", article.ReadAmount)
	return article.ReadAmount, nil
}

// UpdateArticleReadNumById 通过文章id修改文章浏览量
func UpdateArticleReadNumById(aid, num int, db *gorm.DB) error {
	if err := db.Model(&gorm_model.Article{}).Where("id = ?", aid).Update("read_amount", num).Error; err != nil {
		zap.L().Error("UpdateArticleReadNumById() dao.mysql.sql_article", zap.Error(err))
		return err
	}
	return nil
}
