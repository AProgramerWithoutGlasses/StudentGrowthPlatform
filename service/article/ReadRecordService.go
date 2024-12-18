package article

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"studentGrow/dao/mysql"
)

// UpdateArticleReadNumService 增加或减少文章阅读量
func UpdateArticleReadNumService(aid, num int, db *gorm.DB) error {
	// 获取当前阅读量
	curNum, err := mysql.QueryArticleReadNumById(aid)
	if err != nil {
		zap.L().Error("UpdateArticleReadNumService() dao.mysql.sql_article", zap.Error(err))
		return err
	}
	// 更新阅读量
	err = mysql.UpdateArticleReadNumById(aid, curNum+num, db)
	if err != nil {
		zap.L().Error("UpdateArticleReadNumService() dao.mysql.sql_article.UpdateArticleReadNumById", zap.Error(err))
		return err
	}
	return nil
}
